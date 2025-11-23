package handlers

import (
	"lomi-backend/internal/database"
	"lomi-backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// SubmitVerification submits ID verification documents
func SubmitVerification(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	userID, _ := uuid.Parse(userIDStr)

	var req struct {
		SelfieURL     string `json:"selfie_url"`
		IDDocumentURL string `json:"id_document_url"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Check if user already has a pending or approved verification
	var existing models.Verification
	if err := database.DB.Where("user_id = ? AND status IN ?", userID, []string{"pending", "approved"}).
		First(&existing).Error; err == nil {
		if existing.Status == models.VerificationApproved {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User already verified"})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Verification already pending"})
	}

	verification := models.Verification{
		UserID:         userID,
		SelfieURL:      req.SelfieURL,
		IDDocumentURL:  req.IDDocumentURL,
		Status:         models.VerificationPending,
	}

	if err := database.DB.Create(&verification).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to submit verification"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":      "Verification submitted successfully",
		"verification": verification,
	})
}

// GetVerificationStatus returns the current verification status
func GetVerificationStatus(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	userID, _ := uuid.Parse(userIDStr)

	var verification models.Verification
	if err := database.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		First(&verification).Error; err != nil {
		return c.JSON(fiber.Map{
			"is_verified": false,
			"status":      "none",
		})
	}

	return c.JSON(fiber.Map{
		"is_verified": verification.Status == models.VerificationApproved,
		"status":      verification.Status,
		"verification": verification,
	})
}

