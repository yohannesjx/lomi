package handlers

import (
	"fmt"
	"lomi-backend/internal/database"
	"lomi-backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// GetCoinBalance returns the current user's coin balance
func GetCoinBalance(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	userID, _ := uuid.Parse(userIDStr)

	var user models.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(fiber.Map{
		"coin_balance": user.CoinBalance,
		"gift_balance": user.GiftBalance,
	})
}

// PurchaseCoins initiates a coin purchase (creates pending transaction)
func PurchaseCoins(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	userID, _ := uuid.Parse(userIDStr)

	var req struct {
		CoinAmount    int    `json:"coin_amount"`
		PaymentMethod string `json:"payment_method"` // telebirr, cbe_birr, hellocash, amole
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if req.CoinAmount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid coin amount"})
	}

	// Calculate Birr amount (1 coin = 0.10 Birr, or configurable)
	coinToBirrRate := 0.10
	birrAmount := float64(req.CoinAmount) * coinToBirrRate

	// Create pending transaction
	transaction := models.CoinTransaction{
		UserID:          userID,
		TransactionType: models.TransactionTypePurchase,
		CoinAmount:      req.CoinAmount,
		BirrAmount:      birrAmount,
		PaymentMethod:   models.PaymentMethod(req.PaymentMethod),
		PaymentStatus:   models.PaymentStatusPending,
		BalanceAfter:    0, // Will be updated after payment confirmation
	}

	if err := database.DB.Create(&transaction).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create transaction"})
	}

	// Generate payment URL based on payment method
	paymentURL := generatePaymentURL(models.PaymentMethod(req.PaymentMethod), transaction.ID.String(), birrAmount, userID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"transaction_id": transaction.ID,
		"coin_amount":    req.CoinAmount,
		"birr_amount":    birrAmount,
		"payment_method": req.PaymentMethod,
		"payment_url":    paymentURL,
		"status":         "pending",
	})
}

// generatePaymentURL generates payment gateway redirect URL
func generatePaymentURL(paymentMethod models.PaymentMethod, transactionID string, amount float64, userID uuid.UUID) string {
	// Base URL for payment gateways (configure these in your config)
	baseURL := "https://payment.lomi.app" // Replace with actual payment gateway URLs

	switch paymentMethod {
	case models.PaymentMethodTelebirr:
		// Telebirr payment URL format
		// In production, integrate with Telebirr API
		return fmt.Sprintf("%s/telebirr/pay?transaction_id=%s&amount=%.2f&user_id=%s", baseURL, transactionID, amount, userID)
	case models.PaymentMethodCbeBirr:
		// CBE Birr payment URL format
		return fmt.Sprintf("%s/cbe-birr/pay?transaction_id=%s&amount=%.2f&user_id=%s", baseURL, transactionID, amount, userID)
	case models.PaymentMethodHelloCash:
		// HelloCash payment URL format
		return fmt.Sprintf("%s/hellocash/pay?transaction_id=%s&amount=%.2f&user_id=%s", baseURL, transactionID, amount, userID)
	case models.PaymentMethodAmole:
		// Amole payment URL format
		return fmt.Sprintf("%s/amole/pay?transaction_id=%s&amount=%.2f&user_id=%s", baseURL, transactionID, amount, userID)
	default:
		return ""
	}
}

// ConfirmCoinPurchase confirms a coin purchase (called by payment gateway webhook)
func ConfirmCoinPurchase(c *fiber.Ctx) error {
	var req struct {
		TransactionID string `json:"transaction_id"`
		PaymentReference string `json:"payment_reference"`
		Status         string `json:"status"` // completed, failed
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	transactionID, err := uuid.Parse(req.TransactionID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid transaction ID"})
	}

	var transaction models.CoinTransaction
	if err := database.DB.First(&transaction, "id = ?", transactionID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Transaction not found"})
	}

	if transaction.PaymentStatus != models.PaymentStatusPending {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Transaction already processed"})
	}

	tx := database.DB.Begin()

	if req.Status == "completed" {
		transaction.PaymentStatus = models.PaymentStatusCompleted
		transaction.PaymentReference = req.PaymentReference

		// Add coins to user
		var user models.User
		if err := tx.First(&user, "id = ?", transaction.UserID).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}

		user.CoinBalance += transaction.CoinAmount
		transaction.BalanceAfter = user.CoinBalance

		if err := tx.Save(&user).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update balance"})
		}
	} else {
		transaction.PaymentStatus = models.PaymentStatusFailed
	}

	if err := tx.Save(&transaction).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update transaction"})
	}

	tx.Commit()

	return c.JSON(fiber.Map{
		"message": "Transaction updated",
		"status":  transaction.PaymentStatus,
	})
}

// GetCoinTransactions returns transaction history
func GetCoinTransactions(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	userID, _ := uuid.Parse(userIDStr)

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 50)
	offset := (page - 1) * limit

	var transactions []models.CoinTransaction
	if err := database.DB.Where("user_id = ?", userID).
		Preload("GiftTransaction").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch transactions"})
	}

	return c.JSON(fiber.Map{
		"transactions": transactions,
		"page":         page,
		"limit":        limit,
	})
}

