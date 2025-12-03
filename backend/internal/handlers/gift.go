package handlers

import (
	"lomi-backend/internal/database"
	"lomi-backend/internal/models"
	"lomi-backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// GetGifts returns the gift catalog
func GetGifts(c *fiber.Ctx) error {
	var gifts []models.Gift
	if err := database.DB.Where("is_active = ?", true).
		Order("display_order ASC, created_at DESC").
		Find(&gifts).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch gifts"})
	}

	return c.JSON(fiber.Map{
		"gifts": gifts,
		"count": len(gifts),
	})
}

// SendGift sends a gift to a user
func SendGift(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	senderID, _ := uuid.Parse(userIDStr)

	var req struct {
		ReceiverID string `json:"receiver_id"`
		GiftID     string `json:"gift_id"`
		MatchID    string `json:"match_id,omitempty"` // Optional: if sent in chat
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	receiverID, err := uuid.Parse(req.ReceiverID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid receiver ID"})
	}

	giftID, err := uuid.Parse(req.GiftID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid gift ID"})
	}

	// Get gift details
	var gift models.Gift
	if err := database.DB.First(&gift, "id = ? AND is_active = ?", giftID, true).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Gift not found"})
	}

	// Get sender
	var sender models.User
	if err := database.DB.First(&sender, "id = ?", senderID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Sender not found"})
	}

	// Check if sender has enough coins
	if sender.CoinBalance < gift.CoinPrice {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Insufficient coins"})
	}

	// Get receiver
	var receiver models.User
	if err := database.DB.First(&receiver, "id = ?", receiverID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Receiver not found"})
	}

	// Start transaction
	tx := database.DB.Begin()

	// Deduct coins from sender
	sender.CoinBalance -= gift.CoinPrice
	if err := tx.Save(&sender).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to deduct coins"})
	}

	// Add gift value to receiver's gift balance
	receiver.GiftBalance += gift.BirrValue
	if err := tx.Save(&receiver).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add gift balance"})
	}

	// Create gift transaction
	giftTransaction := models.GiftTransaction{
		SenderID:   senderID,
		ReceiverID: receiverID,
		GiftID:     giftID,
		CoinAmount: gift.CoinPrice,
		BirrValue:  gift.BirrValue,
	}
	if err := tx.Create(&giftTransaction).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create gift transaction"})
	}

	// Create coin transaction for sender
	coinTx := models.CoinTransaction{
		UserID:            senderID,
		TransactionType:   models.TransactionTypeGiftSent,
		CoinAmount:        -gift.CoinPrice,
		BalanceAfter:      sender.CoinBalance,
		GiftTransactionID: &giftTransaction.ID,
	}
	if err := tx.Create(&coinTx).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create coin transaction"})
	}

	// Create coin transaction for receiver (gift received)
	coinTxReceiver := models.CoinTransaction{
		UserID:            receiverID,
		TransactionType:   models.TransactionTypeGiftReceived,
		CoinAmount:        0, // No coins, but gift value added to balance
		BalanceAfter:      receiver.CoinBalance,
		GiftTransactionID: &giftTransaction.ID,
	}
	if err := tx.Create(&coinTxReceiver).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create receiver coin transaction"})
	}

	// If sent in chat, create a message
	if req.MatchID != "" {
		matchID, _ := uuid.Parse(req.MatchID)
		message := models.Message{
			MatchID:     &matchID,
			SenderID:    senderID,
			ReceiverID:  &receiverID,
			MessageType: models.MessageTypeGift,
			GiftID:      &giftID,
			IsLive:      false,
		}
		if err := tx.Create(&message).Error; err == nil {
			giftTransaction.MessageID = &message.ID
			tx.Save(&giftTransaction)
		}
	}

	tx.Commit()

	// Send push notification (async)
	go func() {
		if services.NotificationSvc != nil {
			services.NotificationSvc.NotifyGiftReceived(giftTransaction, sender)
		}
	}()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":          "Gift sent successfully",
		"gift_transaction": giftTransaction,
		"gift":             gift,
	})
}
