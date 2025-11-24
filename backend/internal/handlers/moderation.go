package handlers

import (
	"context"
	"fmt"
	"log"
	"lomi-backend/config"
	"lomi-backend/internal/database"
	"lomi-backend/internal/models"
	"lomi-backend/internal/queue"
	"lomi-backend/internal/utils"
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

	// Create media records and prepare job photos
	jobPhotos := make([]queue.PhotoJob, 0, len(req.Photos))
	mediaRecords := make([]models.Media, 0, len(req.Photos))

	for i, photo := range req.Photos {
		// Validate media type
		if photo.MediaType != "photo" && photo.MediaType != "video" {
			continue // Skip invalid types
		}

		// Create media record
		media := models.Media{
			UserID:           userID,
			MediaType:        models.MediaType(photo.MediaType),
			URL:              photo.FileKey, // Store S3 key
			DisplayOrder:     i,
			IsApproved:       false,
			ModerationStatus: "pending",
			BatchID:          batchID,
		}

		if err := database.DB.Create(&media).Error; err != nil {
			log.Printf("❌ Failed to create media record: %v", err)
			continue // Skip failed records
		}

		// Determine bucket based on media type
		var bucket string
		if photo.MediaType == "photo" {
			bucket = config.Cfg.S3BucketPhotos
		} else {
			bucket = config.Cfg.S3BucketVideos
		}

		// Generate presigned download URL for worker (valid for 1 hour)
		ctx := c.Context()
		r2URL, err := database.GeneratePresignedDownloadURL(ctx, bucket, photo.FileKey, 1*time.Hour)
		if err != nil {
			log.Printf("❌ Failed to generate presigned download URL for %s: %v", photo.FileKey, err)
			// Fallback: construct public URL (if bucket is public)
			r2URL = fmt.Sprintf("%s/%s/%s",
				config.Cfg.S3Endpoint,
				bucket,
				photo.FileKey,
			)
			log.Printf("⚠️ Using fallback public URL: %s", r2URL)
		} else {
			log.Printf("✅ Generated presigned download URL for %s (expires in 1h)", photo.FileKey)
		}

		jobPhotos = append(jobPhotos, queue.PhotoJob{
			MediaID: media.ID.String(),
			R2URL:   r2URL,
			R2Key:   photo.FileKey,
			Bucket:  bucket,
		})

		mediaRecords = append(mediaRecords, media)
	}

	if len(jobPhotos) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No valid photos to moderate",
		})
	}

	// Increment rate limit counter
	pipe := database.RedisClient.Pipeline()
	pipe.Incr(ctx, rateLimitKey)
	pipe.Expire(ctx, rateLimitKey, 24*time.Hour)
	pipe.Exec(ctx)

	// Enqueue moderation job (one job for entire batch)
	telegramID := utils.TelegramIDValue(dbUser.TelegramID)
	if err := queue.EnqueuePhotoModeration(batchID, userID, telegramID, jobPhotos); err != nil {
		log.Printf("❌ Failed to enqueue moderation job: %v", err)
		// Don't fail the request - photos are saved, moderation will retry
	}

	log.Printf("✅ Upload complete: batch_id=%s, user_id=%s, photos=%d",
		batchID, userID, len(jobPhotos))

	// Return immediate response (user doesn't wait)
	return c.JSON(fiber.Map{
		"batch_id":     batchID.String(),
		"message":      "We'll check your photos now",
		"photos_count": len(jobPhotos),
		"status":       "pending",
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
