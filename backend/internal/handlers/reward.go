package handlers

import (
	"lomi-backend/internal/database"
	"lomi-backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// GetRewardChannels returns list of reward channels
func GetRewardChannels(c *fiber.Ctx) error {
	var channels []models.RewardChannel
	if err := database.DB.Where("is_active = ?", true).
		Order("display_order ASC, created_at DESC").
		Find(&channels).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch channels"})
	}

	// Check which channels user has already claimed
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	userID, _ := uuid.Parse(userIDStr)

	var claimedChannelIDs []uuid.UUID
	database.DB.Model(&models.UserChannelReward{}).
		Where("user_id = ?", userID).
		Pluck("channel_id", &claimedChannelIDs)

	type ChannelResponse struct {
		models.RewardChannel
		IsClaimed bool `json:"is_claimed"`
	}

	response := make([]ChannelResponse, 0)
	for _, channel := range channels {
		isClaimed := false
		for _, claimedID := range claimedChannelIDs {
			if channel.ID == claimedID {
				isClaimed = true
				break
			}
		}
		response = append(response, ChannelResponse{
			RewardChannel: channel,
			IsClaimed:     isClaimed,
		})
	}

	return c.JSON(fiber.Map{
		"channels": response,
		"count":    len(response),
	})
}

// ClaimChannelReward claims coins for subscribing to a Telegram channel
func ClaimChannelReward(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	userID, _ := uuid.Parse(userIDStr)

	var req struct {
		ChannelID string `json:"channel_id"`
		// In production, you would verify subscription via Telegram Bot API
		// For now, we'll trust the client (not recommended for production)
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	channelID, err := uuid.Parse(req.ChannelID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid channel ID"})
	}

	// Get channel
	var channel models.RewardChannel
	if err := database.DB.First(&channel, "id = ? AND is_active = ?", channelID, true).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Channel not found"})
	}

	// Check if already claimed
	var existing models.UserChannelReward
	if err := database.DB.Where("user_id = ? AND channel_id = ?", userID, channelID).First(&existing).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Reward already claimed"})
	}

	// TODO: Verify user is actually subscribed to the channel via Telegram Bot API
	// For now, we'll proceed without verification

	tx := database.DB.Begin()

	// Create reward record
	reward := models.UserChannelReward{
		UserID:       userID,
		ChannelID:    channelID,
		RewardAmount: channel.CoinReward,
	}
	if err := tx.Create(&reward).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create reward"})
	}

	// Add coins to user
	var user models.User
	if err := tx.First(&user, "id = ?", userID).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	user.CoinBalance += channel.CoinReward
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update balance"})
	}

	// Create coin transaction
	coinTx := models.CoinTransaction{
		UserID:          userID,
		TransactionType: models.TransactionTypeChannelSubscriptionReward,
		CoinAmount:      channel.CoinReward,
		BalanceAfter:    user.CoinBalance,
		Metadata:        models.JSONMap{"channel_id": channelID.String()},
	}
	if err := tx.Create(&coinTx).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create transaction"})
	}

	tx.Commit()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":      "Reward claimed successfully",
		"coins_earned": channel.CoinReward,
		"new_balance":  user.CoinBalance,
	})
}

