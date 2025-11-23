package handlers

import (
	"lomi-backend/internal/database"
	"lomi-backend/internal/models"
	"lomi-backend/internal/services"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// GetChats returns all conversations for the current user
func GetChats(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	userID, _ := uuid.Parse(userIDStr)

	var matches []models.Match
	if err := database.DB.Where("(user1_id = ? OR user2_id = ?) AND is_active = ?", userID, userID, true).
		Preload("User1").
		Preload("User2").
		Find(&matches).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch chats"})
	}

	type ChatResponse struct {
		MatchID     uuid.UUID       `json:"match_id"`
		User        models.User      `json:"user"`
		LastMessage *models.Message  `json:"last_message,omitempty"`
		UnreadCount int64            `json:"unread_count"`
	}

	chats := make([]ChatResponse, 0)
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

		// Count unread messages
		var unreadCount int64
		database.DB.Model(&models.Message{}).
			Where("match_id = ? AND receiver_id = ? AND is_read = ?", match.ID, userID, false).
			Count(&unreadCount)

		chats = append(chats, ChatResponse{
			MatchID:     match.ID,
			User:        otherUser,
			LastMessage: &lastMessage,
			UnreadCount: unreadCount,
		})
	}

	return c.JSON(fiber.Map{
		"chats": chats,
		"count": len(chats),
	})
}

// GetMessages returns messages for a specific match
func GetMessages(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	userID, _ := uuid.Parse(userIDStr)

	matchID := c.Params("id")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 50)
	offset := (page - 1) * limit

	// Verify user is part of this match
	var match models.Match
	if err := database.DB.Where("id = ? AND (user1_id = ? OR user2_id = ?) AND is_active = ?", matchID, userID, userID, true).
		First(&match).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Match not found"})
	}

	var messages []models.Message
	if err := database.DB.Where("match_id = ?", matchID).
		Preload("Sender").
		Preload("Receiver").
		Preload("Gift").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch messages"})
	}

	// Mark messages as read
	now := time.Now()
	database.DB.Model(&models.Message{}).
		Where("match_id = ? AND receiver_id = ? AND is_read = ?", matchID, userID, false).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		})

	return c.JSON(fiber.Map{
		"messages": messages,
		"page":     page,
		"limit":    limit,
	})
}

// SendMessage creates a new message
func SendMessage(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	senderID, _ := uuid.Parse(userIDStr)

	var req struct {
		MatchID     string                 `json:"match_id"`
		MessageType string                 `json:"message_type"`
		Content     string                 `json:"content,omitempty"`
		MediaURL    string                 `json:"media_url,omitempty"`
		GiftID      string                 `json:"gift_id,omitempty"`
		Metadata    map[string]interface{} `json:"metadata,omitempty"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	matchID, err := uuid.Parse(req.MatchID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid match ID"})
	}

	// Verify user is part of this match
	var match models.Match
	if err := database.DB.Where("id = ? AND (user1_id = ? OR user2_id = ?) AND is_active = ?", matchID, senderID, senderID, true).
		First(&match).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Match not found"})
	}

	// Determine receiver
	var receiverID uuid.UUID
	if match.User1ID == senderID {
		receiverID = match.User2ID
	} else {
		receiverID = match.User1ID
	}

	// Check if sender is blocked by receiver
	var block models.Block
	if err := database.DB.Where("(blocker_id = ? AND blocked_id = ?) OR (blocker_id = ? AND blocked_id = ?)", receiverID, senderID, senderID, receiverID).First(&block).Error; err == nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Cannot send message: user is blocked"})
	}

	// Create message
	message := models.Message{
		MatchID:     matchID,
		SenderID:    senderID,
		ReceiverID:  receiverID,
		MessageType: models.MessageType(req.MessageType),
		Content:     req.Content,
		MediaURL:    req.MediaURL,
		IsRead:      false,
	}

	if req.GiftID != "" {
		giftID, _ := uuid.Parse(req.GiftID)
		message.GiftID = &giftID
	}

	if req.Metadata != nil {
		message.Metadata = models.JSONMap(req.Metadata)
	}

	if err := database.DB.Create(&message).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to send message"})
	}

	// Load relations for response
	database.DB.Preload("Sender").Preload("Receiver").Preload("Gift").First(&message, message.ID)

	// Send push notification (async)
	go func() {
		if services.NotificationSvc != nil {
			services.NotificationSvc.NotifyNewMessage(message, message.Sender)
		}
	}()

	return c.Status(fiber.StatusCreated).JSON(message)
}

// MarkMessagesAsRead marks all messages in a match as read
func MarkMessagesAsRead(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	userID, _ := uuid.Parse(userIDStr)

	matchID := c.Params("id")

	// Verify user is part of this match
	var match models.Match
	if err := database.DB.Where("id = ? AND (user1_id = ? OR user2_id = ?) AND is_active = ?", matchID, userID, userID, true).
		First(&match).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Match not found"})
	}

	now := time.Now()
	if err := database.DB.Model(&models.Message{}).
		Where("match_id = ? AND receiver_id = ? AND is_read = ?", matchID, userID, false).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to mark as read"})
	}

	return c.JSON(fiber.Map{"message": "Messages marked as read"})
}

