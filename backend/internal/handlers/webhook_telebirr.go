package handlers

import (
	"log"
	"lomi-backend/internal/database"
	"lomi-backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CoinPurchaseWebhook handles Telebirr/CBE Birr payment webhook
func CoinPurchaseWebhook(c *fiber.Ctx) error {
	// TODO: Verify webhook signature from Telebirr
	// For now, accept the webhook data

	var webhookData struct {
		TransactionID string  `json:"transaction_id"` // Our internal transaction ID
		PaymentRef    string  `json:"payment_reference"`
		Amount        float64 `json:"amount"`
		Status        string  `json:"status"` // "success", "failed", "pending"
		PhoneNumber   string  `json:"phone_number"`
	}

	if err := c.BodyParser(&webhookData); err != nil {
		log.Printf("❌ Webhook parse error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid webhook data"})
	}

	// Find transaction
	transactionID, err := uuid.Parse(webhookData.TransactionID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid transaction ID"})
	}

	var coinTx models.CoinTransaction
	if err := database.DB.First(&coinTx, "id = ?", transactionID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Transaction not found"})
	}

	// Check if already processed
	if coinTx.PaymentStatus == models.PaymentStatusCompleted {
		return c.JSON(fiber.Map{"message": "Transaction already processed"})
	}

	// Process payment
	if webhookData.Status == "success" {
		// Update transaction status
		coinTx.PaymentStatus = models.PaymentStatusCompleted
		coinTx.PaymentReference = webhookData.PaymentRef
		if err := database.DB.Save(&coinTx).Error; err != nil {
			log.Printf("❌ Failed to update transaction: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update transaction"})
		}

		// Add coins to user
		var user models.User
		if err := database.DB.First(&user, "id = ?", coinTx.UserID).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}

		user.CoinBalance += coinTx.CoinAmount
		coinTx.BalanceAfter = user.CoinBalance
		if err := database.DB.Save(&user).Error; err != nil {
			log.Printf("❌ Failed to add coins to user: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add coins"})
		}

		// Update balance_after in transaction
		database.DB.Save(&coinTx)

		log.Printf("✅ Payment successful: User %s received %d coins", user.ID, coinTx.CoinAmount)
		return c.JSON(fiber.Map{
			"message": "Payment processed successfully",
			"coins_added": coinTx.CoinAmount,
			"new_balance": user.CoinBalance,
		})
	} else if webhookData.Status == "failed" {
		coinTx.PaymentStatus = models.PaymentStatusFailed
		database.DB.Save(&coinTx)
		return c.JSON(fiber.Map{"message": "Payment failed"})
	}

	return c.JSON(fiber.Map{"message": "Webhook received"})
}

