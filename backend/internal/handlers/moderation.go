package handlers

import (
	"context"
	"fmt"
	"log"
	"lomi-backend/config"
	"lomi-backend/internal/database"
	"lomi-backend/internal/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// UploadComplete handles batch photo upload completion and enqueues moderation
func UploadComplete(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	userID, _ := uuid.Parse(userIDStr)

	var req struct {
		Photos []struct {
			FileKey   string `json:"file_key"`
			MediaType string `json:"media_type"`
		} `json:"photos"`
	}

	if err := c.BodyParser(&req); err != nil {
		log.Printf("❌ Failed to parse upload-complete request: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	// Validate batch size (1-9 photos)
	if len(req.Photos) == 0 || len(req.Photos) > 9 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Must upload between 1 and 9 photos",
		})
	}

	// Rate limit check: 30 photos per 24 hours
	rateLimitKey := fmt.Sprintf("photo_upload_rate:%s", userID.String())
	ctx := c.Context()

	currentCount, err := database.RedisClient.Get(ctx, rateLimitKey).Int()
	if err != nil && err.Error() != "redis: nil" {
		log.Printf("⚠️ Rate limit check failed: %v", err)
		// Continue on error (don't block user)
	} else if currentCount >= 30 {
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"error":       "Rate limit exceeded",
			"message":     "Maximum 30 photos per 24 hours. Please try again tomorrow.",
			"retry_after": 86400,
		})
	}

	// Get user's Telegram ID for push notifications
	var dbUser models.User
	if err := database.DB.First(&dbUser, "id = ?", userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Generate batch_id for this upload session
	batchID := uuid.New()

	// Create media records and auto-approve them (manual moderation)
	mediaRecords := make([]models.Media, 0, len(req.Photos))
	approvedCount := 0

	for i, photo := range req.Photos {
		// Validate media type
		if photo.MediaType != "photo" && photo.MediaType != "video" {
			continue // Skip invalid types
		}

		// Create media record - auto-approve (manual moderation will be done later)
		media := models.Media{
			UserID:           userID,
			MediaType:        models.MediaType(photo.MediaType),
			URL:              photo.FileKey, // Store S3 key
			DisplayOrder:     i,
			IsApproved:       true, // Auto-approve for now
			ModerationStatus: "approved", // Auto-approve
			BatchID:          batchID,
		}

		if err := database.DB.Create(&media).Error; err != nil {
			log.Printf("❌ Failed to create media record: %v", err)
			continue // Skip failed records
		}

		mediaRecords = append(mediaRecords, media)
		approvedCount++
	}

	if len(mediaRecords) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No valid photos to upload",
		})
	}

	// Increment rate limit counter
	pipe := database.RedisClient.Pipeline()
	pipe.Incr(ctx, rateLimitKey)
	pipe.Expire(ctx, rateLimitKey, 24*time.Hour)
	pipe.Exec(ctx)

	log.Printf("✅ Upload complete: batch_id=%s, user_id=%s, photos=%d (auto-approved)",
		batchID, userID, approvedCount)

	// Return immediate response
	return c.JSON(fiber.Map{
		"batch_id":     batchID.String(),
		"message":      "Photos uploaded successfully",
		"photos_count": approvedCount,
		"status":       "approved",
	})
}

// GetMyModerationStatus returns moderation status for the authenticated user's media
func GetMyModerationStatus(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var media []models.Media
	if err := database.DB.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&media).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to fetch media",
			"details": err.Error(),
		})
	}

	ctx := context.Background()
	expiresIn := 12 * time.Hour

	type summaryCounters struct {
		Total    int
		Approved int
		Pending  int
		Rejected int
		Failed   int
	}

	summary := summaryCounters{}
	lastModeratedAt := time.Time{}
	pendingBatchMap := make(map[uuid.UUID]struct{})

	photos := make([]fiber.Map, 0, len(media))

	for _, m := range media {
		summary.Total++
		switch m.ModerationStatus {
		case "approved":
			summary.Approved++
		case "rejected":
			summary.Rejected++
		case "failed":
			summary.Failed++
		default:
			summary.Pending++
			if m.BatchID != uuid.Nil {
				pendingBatchMap[m.BatchID] = struct{}{}
			}
		}

		if !m.ModeratedAt.IsZero() && m.ModeratedAt.After(lastModeratedAt) {
			lastModeratedAt = m.ModeratedAt
		}

		bucket := config.Cfg.S3BucketPhotos
		if m.MediaType == models.MediaTypeVideo {
			bucket = config.Cfg.S3BucketVideos
		}

		downloadURL := ""
		if m.URL != "" {
			if url, err := database.GeneratePresignedDownloadURL(ctx, bucket, m.URL, expiresIn); err == nil {
				downloadURL = url
			}
		}

		photos = append(photos, fiber.Map{
			"id":                m.ID,
			"batch_id":          m.BatchID,
			"media_type":        m.MediaType,
			"status":            m.ModerationStatus,
			"reason":            m.ModerationReason,
			"moderated_at":      m.ModeratedAt,
			"uploaded_at":       m.CreatedAt,
			"display_order":     m.DisplayOrder,
			"url":               downloadURL,
			"scores":            m.ModerationScores,
			"retry_count":       m.RetryCount,
			"is_approved":       m.IsApproved,
			"moderation_status": m.ModerationStatus,
		})
	}

	needsMorePhotos := summary.Approved < 2
	pendingBatchCount := len(pendingBatchMap)

	response := fiber.Map{
		"summary": fiber.Map{
			"total_photos":      summary.Total,
			"approved":          summary.Approved,
			"pending":           summary.Pending,
			"rejected":          summary.Rejected,
			"failed":            summary.Failed,
			"needs_more_photos": needsMorePhotos,
			"pending_batches":   pendingBatchCount,
			"last_moderated_at": lastModeratedAt,
		},
		"photos": photos,
	}

	return c.JSON(response)
}
