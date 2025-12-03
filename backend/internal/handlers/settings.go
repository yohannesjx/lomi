package handlers

import (
	"lomi-backend/internal/database"
	"lomi-backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// GetPrivacySettings retrieves the user's privacy settings
func GetPrivacySettings(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	var settings models.PrivacySetting
	if err := database.DB.Where("user_id = ?", userID).First(&settings).Error; err != nil {
		// If not found, create default settings
		settings = models.PrivacySetting{
			UserID: uuid.MustParse(userID),
		}
		database.DB.Create(&settings)
	}

	return c.JSON(settings)
}

// UpdatePrivacySettings updates the user's privacy settings
func UpdatePrivacySettings(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	var req models.PrivacySetting
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var settings models.PrivacySetting
	if err := database.DB.Where("user_id = ?", userID).First(&settings).Error; err != nil {
		// If not found, create new
		settings = models.PrivacySetting{
			UserID: uuid.MustParse(userID),
		}
	}

	// Update fields
	settings.VideosDownload = req.VideosDownload
	settings.DirectMessage = req.DirectMessage
	settings.Duet = req.Duet
	settings.LikedVideos = req.LikedVideos
	settings.VideoComment = req.VideoComment
	settings.OrderHistory = req.OrderHistory

	if err := database.DB.Save(&settings).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update privacy settings"})
	}

	return c.JSON(settings)
}

// GetPushNotifications retrieves the user's push notification settings
func GetPushNotifications(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	var settings models.PushNotification
	if err := database.DB.Where("user_id = ?", userID).First(&settings).Error; err != nil {
		// If not found, create default settings
		settings = models.PushNotification{
			UserID: uuid.MustParse(userID),
		}
		database.DB.Create(&settings)
	}

	return c.JSON(settings)
}

// UpdatePushNotifications updates the user's push notification settings
func UpdatePushNotifications(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	var req models.PushNotification
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var settings models.PushNotification
	if err := database.DB.Where("user_id = ?", userID).First(&settings).Error; err != nil {
		// If not found, create new
		settings = models.PushNotification{
			UserID: uuid.MustParse(userID),
		}
	}

	// Update fields
	settings.Likes = req.Likes
	settings.Comments = req.Comments
	settings.NewFollowers = req.NewFollowers
	settings.Mentions = req.Mentions
	settings.DirectMessages = req.DirectMessages
	settings.VideoUpdates = req.VideoUpdates

	if err := database.DB.Save(&settings).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update push notification settings"})
	}

	return c.JSON(settings)
}
