package handlers

import (
	"context"
	"log"
	"lomi-backend/config"
	"lomi-backend/internal/database"
	"lomi-backend/internal/models"
	"lomi-backend/internal/services"
	"math"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// GetSwipeCards returns potential matches for swiping
func GetSwipeCards(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	userID, _ := uuid.Parse(userIDStr)

	var currentUser models.User
	if err := database.DB.First(&currentUser, "id = ?", userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Get user preferences
	var minAge, maxAge int = 18, 100
	// TODO: Implement distance-based filtering when location data is available
	// var maxDistance float64 = 50 // km
	if prefs, ok := currentUser.Preferences["min_age"].(float64); ok {
		minAge = int(prefs)
	}
	if prefs, ok := currentUser.Preferences["max_age"].(float64); ok {
		maxAge = int(prefs)
	}
	// if prefs, ok := currentUser.Preferences["max_distance"].(float64); ok {
	// 	maxDistance = prefs
	// }

	// Get already swiped user IDs
	var swipedIDs []uuid.UUID
	database.DB.Model(&models.Swipe{}).
		Where("swiper_id = ?", userID).
		Pluck("swiped_id", &swipedIDs)

	// Get blocked users
	var blockedIDs []uuid.UUID
	database.DB.Model(&models.Block{}).
		Where("blocker_id = ? OR blocked_id = ?", userID, userID).
		Pluck("blocked_id", &swipedIDs)
	database.DB.Model(&models.Block{}).
		Where("blocker_id = ? OR blocked_id = ?", userID, userID).
		Pluck("blocker_id", &swipedIDs)

	// Build query
	query := database.DB.Where("id != ?", userID).
		Where("is_active = ?", true).
		Where("age >= ? AND age <= ?", minAge, maxAge)
		// Temporarily disabled city filter for testing
		// .Where("city = ?", currentUser.City)

	// Only exclude if there are actually blocked/swiped users
	if len(swipedIDs) > 0 {
		query = query.Where("id NOT IN ?", swipedIDs)
	}
	if len(blockedIDs) > 0 {
		query = query.Where("id NOT IN ?", blockedIDs)
	}

	// Gender preference - temporarily disabled for testing
	// var lookingFor string
	// if prefs, ok := currentUser.Preferences["looking_for"].(string); ok {
	// 	lookingFor = prefs
	// } else {
	// 	// Default: opposite gender
	// 	if currentUser.Gender == models.GenderMale {
	// 		lookingFor = "female"
	// 	} else if currentUser.Gender == models.GenderFemale {
	// 		lookingFor = "male"
	// 	}
	// }

	// if lookingFor == "male" {
	// 	query = query.Where("gender = ?", models.GenderMale)
	// } else if lookingFor == "female" {
	// 	query = query.Where("gender = ?", models.GenderFemale)
	// }

	// Limit to 20 cards per request
	var users []models.User

	// Debug logging
	log.Printf("🔍 Discover query filters:")
	log.Printf("  - Current user ID: %s", userID)
	log.Printf("  - Age range: %d - %d", minAge, maxAge)
	log.Printf("  - Swiped IDs count: %d", len(swipedIDs))
	log.Printf("  - Blocked IDs count: %d", len(blockedIDs))
	log.Printf("  - is_active: true")

	if err := query.Limit(20).Order("created_at DESC").Find(&users).Error; err != nil {
		log.Printf("❌ Query error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch users"})
	}

	log.Printf("✅ Found %d users", len(users))

	// Get photos for each user
	type UserCard struct {
		User     models.User    `json:"user"`
		Photos   []models.Media `json:"photos"`
		Video    *models.Media  `json:"video,omitempty"`
		Distance float64        `json:"distance"`
	}

	cards := make([]UserCard, 0)
	for _, u := range users {
		var photos []models.Media
		database.DB.Where("user_id = ? AND media_type = ? AND is_approved = ?", u.ID, models.MediaTypePhoto, true).
			Order("display_order ASC").Limit(9).Find(&photos)

		var video models.Media
		hasVideo := database.DB.Where("user_id = ? AND media_type = ? AND is_approved = ?", u.ID, models.MediaTypeVideo, true).
			First(&video).Error == nil

		// Calculate distance (Haversine formula simplified)
		distance := 0.0
		if currentUser.Latitude != 0 && currentUser.Longitude != 0 && u.Latitude != 0 && u.Longitude != 0 {
			distance = calculateDistance(currentUser.Latitude, currentUser.Longitude, u.Latitude, u.Longitude)
		}

		card := UserCard{
			User:     u,
			Photos:   photos,
			Distance: distance,
		}
		if hasVideo {
			card.Video = &video
		}
		cards = append(cards, card)
	}

	return c.JSON(fiber.Map{
		"cards": cards,
		"count": len(cards),
	})
}

// SwipeAction handles like/pass/super_like actions
func SwipeAction(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	swiperID, _ := uuid.Parse(userIDStr)

	var req struct {
		SwipedID string `json:"swiped_id"`
		Action   string `json:"action"` // "like", "pass", "super_like"
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	swipedID, err := uuid.Parse(req.SwipedID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Check if already swiped
	var existingSwipe models.Swipe
	if err := database.DB.Where("swiper_id = ? AND swiped_id = ?", swiperID, swipedID).First(&existingSwipe).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Already swiped"})
	}

	// Create swipe record
	swipe := models.Swipe{
		SwiperID: swiperID,
		SwipedID: swipedID,
		Action:   models.SwipeAction(req.Action),
	}
	if err := database.DB.Create(&swipe).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to record swipe"})
	}

	// Check for match (if action is like or super_like)
	if req.Action == "like" || req.Action == "super_like" {
		var mutualSwipe models.Swipe
		if err := database.DB.Where("swiper_id = ? AND swiped_id = ? AND action IN ?", swipedID, swiperID, []string{"like", "super_like"}).First(&mutualSwipe).Error; err == nil {
			// It's a match!
			match := models.Match{
				User1ID:     swiperID,
				User2ID:     swipedID,
				InitiatedBy: swiperID,
				IsActive:    true,
			}
			// Ensure consistent ordering
			if swiperID.String() > swipedID.String() {
				match.User1ID, match.User2ID = match.User2ID, match.User1ID
			}
			if err := database.DB.Create(&match).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create match"})
			}

			// Get matched user details for notification
			var matchedUser models.User
			database.DB.First(&matchedUser, "id = ?", swipedID)

			// Send push notification (async)
			go func() {
				if services.NotificationSvc != nil {
					services.NotificationSvc.NotifyNewMatch(match, matchedUser)
				}
			}()

			return c.JSON(fiber.Map{
				"match":    true,
				"match_id": match.ID,
				"message":  "It's a match! 💚",
				"user":     matchedUser,
			})
		} else {
			// No match yet, but send "someone liked you" notification if enabled
			var swipedUser models.User
			if err := database.DB.First(&swipedUser, "id = ?", swipedID).Error; err == nil {
				var swiperUser models.User
				if err := database.DB.First(&swiperUser, "id = ?", swiperID).Error; err == nil {
					// Send "someone liked you" notification (optional, can be disabled)
					go func() {
						if services.NotificationSvc != nil {
							services.NotificationSvc.NotifySomeoneLiked(swiperUser, swipedID)
						}
					}()
				}
			}
		}
	}

	return c.JSON(fiber.Map{
		"match":   false,
		"message": "Swipe recorded",
	})
}

// GetExploreFeed returns TikTok-style vertical feed
func GetExploreFeed(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	userID, _ := uuid.Parse(userIDStr)

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)
	offset := (page - 1) * limit

	// Get media from active users (excluding current user)
	var media []models.Media
	query := database.DB.
		Joins("JOIN users ON media.user_id = users.id").
		Where("users.is_active = ?", true).
		Where("users.id != ?", userID).
		Where("media.is_approved = ?", true).
		Order("media.created_at DESC").
		Limit(limit).
		Offset(offset)

	if err := query.Find(&media).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch feed"})
	}

	// Format response with presigned URLs
	type FormattedFeedItem struct {
		Media struct {
			ID           string `json:"id"`
			MediaType    string `json:"media_type"`
			URL          string `json:"url"`
			ThumbnailURL string `json:"thumbnail_url,omitempty"`
		} `json:"media"`
		User struct {
			ID     string `json:"id"`
			Name   string `json:"name"`
			Avatar string `json:"avatar"`
		} `json:"user"`
	}

	ctx := context.Background()
	expiresIn := 24 * time.Hour // URLs valid for 24 hours
	formattedItems := make([]FormattedFeedItem, 0)

	for _, m := range media {
		var u models.User
		database.DB.First(&u, "id = ?", m.UserID)

		// Determine bucket based on media type
		bucket := config.Cfg.S3BucketPhotos
		if m.MediaType == models.MediaTypeVideo {
			bucket = config.Cfg.S3BucketVideos
		}

		// Generate presigned download URL for media
		downloadURL, err := database.GeneratePresignedDownloadURL(ctx, bucket, m.URL, expiresIn)
		if err != nil {
			downloadURL = "" // Fallback if URL generation fails
		}

		// Generate thumbnail URL if exists
		thumbnailURL := ""
		if m.ThumbnailURL != "" {
			thumbnailURL, _ = database.GeneratePresignedDownloadURL(ctx, config.Cfg.S3BucketPhotos, m.ThumbnailURL, expiresIn)
		}

		// Get user's first photo as avatar
		var userPhoto models.Media
		userAvatarURL := ""
		if err := database.DB.Where("user_id = ? AND media_type = ? AND is_approved = ?", u.ID, models.MediaTypePhoto, true).
			Order("display_order ASC").First(&userPhoto).Error; err == nil {
			userAvatarURL, _ = database.GeneratePresignedDownloadURL(ctx, config.Cfg.S3BucketPhotos, userPhoto.URL, expiresIn)
		}

		formattedItems = append(formattedItems, FormattedFeedItem{
			Media: struct {
				ID           string `json:"id"`
				MediaType    string `json:"media_type"`
				URL          string `json:"url"`
				ThumbnailURL string `json:"thumbnail_url,omitempty"`
			}{
				ID:           m.ID.String(),
				MediaType:    string(m.MediaType),
				URL:          downloadURL,
				ThumbnailURL: thumbnailURL,
			},
			User: struct {
				ID     string `json:"id"`
				Name   string `json:"name"`
				Avatar string `json:"avatar"`
			}{
				ID:     u.ID.String(),
				Name:   u.Name,
				Avatar: userAvatarURL,
			},
		})
	}

	return c.JSON(fiber.Map{
		"items": formattedItems,
		"page":  page,
		"limit": limit,
	})
}

// Helper function to calculate distance between two coordinates (Haversine)
func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth radius in km
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}
