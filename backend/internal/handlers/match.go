package handlers

import (
	"lomi-backend/internal/database"
	"lomi-backend/internal/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// GetMatches returns all active matches for the current user
func GetMatches(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	userID, _ := uuid.Parse(userIDStr)

	var matches []models.Match
	if err := database.DB.Where("(user1_id = ? OR user2_id = ?) AND is_active = ?", userID, userID, true).
		Preload("User1").
		Preload("User2").
		Order("created_at DESC").
		Find(&matches).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch matches"})
	}

	// Format response with other user info
	type MatchResponse struct {
		ID        uuid.UUID   `json:"id"`
		User      models.User `json:"user"`
		CreatedAt string      `json:"created_at"`
		LastMessage *models.Message `json:"last_message,omitempty"`
	}

	response := make([]MatchResponse, 0)
	for _, match := range matches {
		var otherUser models.User
		if match.User1ID == userID {
			otherUser = match.User2
		} else {
			otherUser = match.User1
		}

		// Get last message
		var lastMessage models.Message
		database.DB.Where("match_id = ?", match.ID).
			Order("created_at DESC").
			First(&lastMessage)

		response = append(response, MatchResponse{
			ID:        match.ID,
			User:      otherUser,
			CreatedAt: match.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			LastMessage: &lastMessage,
		})
	}

	return c.JSON(fiber.Map{
		"matches": response,
		"count":   len(response),
	})
}

// GetMatchDetails returns details of a specific match
func GetMatchDetails(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	userID, _ := uuid.Parse(userIDStr)

	matchID := c.Params("id")
	var match models.Match
	if err := database.DB.Where("id = ? AND (user1_id = ? OR user2_id = ?) AND is_active = ?", matchID, userID, userID, true).
		Preload("User1").
		Preload("User2").
		First(&match).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Match not found"})
	}

	var otherUser models.User
	if match.User1ID == userID {
		otherUser = match.User2
	} else {
		otherUser = match.User1
	}

	// Get user photos
	var photos []models.Media
	database.DB.Where("user_id = ? AND media_type = ? AND is_approved = ?", otherUser.ID, models.MediaTypePhoto, true).
		Order("display_order ASC").Find(&photos)

	return c.JSON(fiber.Map{
		"match": match,
		"user":  otherUser,
		"photos": photos,
	})
}

// Unmatch removes a match
func Unmatch(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	userID, _ := uuid.Parse(userIDStr)

	matchID := c.Params("id")
	var match models.Match
	if err := database.DB.Where("id = ? AND (user1_id = ? OR user2_id = ?) AND is_active = ?", matchID, userID, userID, true).
		First(&match).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Match not found"})
	}

	now := time.Now()
	match.IsActive = false
	match.UnmatchedBy = &userID
	match.UnmatchedAt = &now

	if err := database.DB.Save(&match).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to unmatch"})
	}

	return c.JSON(fiber.Map{"message": "Unmatched successfully"})
}

