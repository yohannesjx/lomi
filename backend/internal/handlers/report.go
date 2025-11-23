package handlers

import (
	"fmt"
	"lomi-backend/internal/database"
	"lomi-backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// ReportUser reports a user for inappropriate behavior
func ReportUser(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	reporterID, _ := uuid.Parse(userIDStr)

	var req struct {
		ReportedUserID string   `json:"reported_user_id"`
		Reason         string   `json:"reason"`
		Description    string   `json:"description,omitempty"`
		ScreenshotURLs []string `json:"screenshot_urls,omitempty"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	reportedUserID, err := uuid.Parse(req.ReportedUserID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	if reporterID == reportedUserID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot report yourself"})
	}

	report := models.Report{
		ReporterID:      reporterID,
		ReportedUserID:  reportedUserID,
		Reason:          models.ReportReason(req.Reason),
		Description:     req.Description,
		ScreenshotURLs:  models.JSONStringArray(req.ScreenshotURLs),
		IsReviewed:      false,
	}

	if err := database.DB.Create(&report).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create report"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Report submitted successfully",
		"report":  report,
	})
}

// ReportPhoto reports a photo for inappropriate content
func ReportPhoto(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	reporterID, _ := uuid.Parse(userIDStr)

	var req struct {
		MediaID        string   `json:"media_id"`
		Reason         string   `json:"reason"`
		Description    string   `json:"description,omitempty"`
		ScreenshotURLs []string `json:"screenshot_urls,omitempty"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	mediaID, err := uuid.Parse(req.MediaID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid media ID"})
	}

	// Get media to find the owner
	var media models.Media
	if err := database.DB.First(&media, "id = ?", mediaID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Media not found"})
	}

	if reporterID == media.UserID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot report your own photo"})
	}

	// Create report for the photo owner
	report := models.Report{
		ReporterID:      reporterID,
		ReportedUserID:  media.UserID,
		Reason:          models.ReportReason(req.Reason),
		Description:     fmt.Sprintf("Reported photo: %s. %s", mediaID.String(), req.Description),
		ScreenshotURLs:  models.JSONStringArray(req.ScreenshotURLs),
		IsReviewed:      false,
	}

	if err := database.DB.Create(&report).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create report"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Photo report submitted successfully",
		"report":  report,
	})
}

// BlockUser blocks a user
func BlockUser(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	blockerID, _ := uuid.Parse(userIDStr)

	var req struct {
		BlockedUserID string `json:"blocked_user_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	blockedID, err := uuid.Parse(req.BlockedUserID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	if blockerID == blockedID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot block yourself"})
	}

	// Check if already blocked
	var existing models.Block
	if err := database.DB.Where("blocker_id = ? AND blocked_id = ?", blockerID, blockedID).First(&existing).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User already blocked"})
	}

	block := models.Block{
		BlockerID: blockerID,
		BlockedID: blockedID,
	}

	if err := database.DB.Create(&block).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to block user"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User blocked successfully",
	})
}

// UnblockUser unblocks a user
func UnblockUser(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	blockerID, _ := uuid.Parse(userIDStr)

	blockedIDParam := c.Params("user_id")
	blockedID, err := uuid.Parse(blockedIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	if err := database.DB.Where("blocker_id = ? AND blocked_id = ?", blockerID, blockedID).Delete(&models.Block{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to unblock user"})
	}

	return c.JSON(fiber.Map{"message": "User unblocked successfully"})
}

// GetBlockedUsers returns list of blocked users
func GetBlockedUsers(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	blockerID, _ := uuid.Parse(userIDStr)

	var blocks []models.Block
	if err := database.DB.Where("blocker_id = ?", blockerID).
		Preload("Blocked").
		Find(&blocks).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch blocked users"})
	}

	blockedUsers := make([]models.User, 0)
	for _, block := range blocks {
		blockedUsers = append(blockedUsers, block.Blocked)
	}

	return c.JSON(fiber.Map{
		"blocked_users": blockedUsers,
		"count":         len(blockedUsers),
	})
}

