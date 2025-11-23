package handlers

import (
	"lomi-backend/internal/database"
	"lomi-backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// GetPayoutBalance returns the user's available payout balance
func GetPayoutBalance(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	userID, _ := uuid.Parse(userIDStr)

	var user models.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Get pending payout amount
	var pendingAmount float64
	database.DB.Model(&models.Payout{}).
		Where("user_id = ? AND status IN ?", userID, []string{"pending", "processing"}).
		Select("COALESCE(SUM(gift_balance_amount), 0)").
		Scan(&pendingAmount)

	return c.JSON(fiber.Map{
		"available_balance": user.GiftBalance,
		"pending_payouts":   pendingAmount,
		"total_earned":      user.GiftBalance + pendingAmount,
	})
}

// RequestPayout creates a payout request
func RequestPayout(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	userID, _ := uuid.Parse(userIDStr)

	var req struct {
		Amount        float64 `json:"amount"`
		PaymentMethod string   `json:"payment_method"`
		PaymentAccount string  `json:"payment_account"`
		PaymentAccountName string `json:"payment_account_name,omitempty"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	var user models.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Check minimum payout amount (1000 Birr)
	minPayoutAmount := 1000.0
	if req.Amount < minPayoutAmount {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Minimum payout amount is 1000 Birr",
			"minimum": minPayoutAmount,
		})
	}

	// Check if user has enough balance
	if user.GiftBalance < req.Amount {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Insufficient balance"})
	}

	// Calculate platform fee (25%)
	platformFeePercentage := 25
	platformFeeAmount := req.Amount * float64(platformFeePercentage) / 100.0
	netAmount := req.Amount - platformFeeAmount

	// Create payout request
	payout := models.Payout{
		UserID:              userID,
		GiftBalanceAmount:   req.Amount,
		PlatformFeePercentage: platformFeePercentage,
		PlatformFeeAmount:   platformFeeAmount,
		NetAmount:           netAmount,
		PaymentMethod:       models.PaymentMethod(req.PaymentMethod),
		PaymentAccount:      req.PaymentAccount,
		PaymentAccountName:  req.PaymentAccountName,
		Status:              models.PayoutStatusPending,
	}

	if err := database.DB.Create(&payout).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create payout request"})
	}

	// Deduct from user's gift balance
	user.GiftBalance -= req.Amount
	if err := database.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update balance"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Payout request created",
		"payout":  payout,
	})
}

// GetPayoutHistory returns payout history
func GetPayoutHistory(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	userID, _ := uuid.Parse(userIDStr)

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 50)
	offset := (page - 1) * limit

	var payouts []models.Payout
	if err := database.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&payouts).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch payouts"})
	}

	return c.JSON(fiber.Map{
		"payouts": payouts,
		"page":    page,
		"limit":   limit,
	})
}

