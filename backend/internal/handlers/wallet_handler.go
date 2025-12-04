package handlers

import (
	"strconv"

	"lomi-backend/internal/models"
	"lomi-backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// WalletHandler handles wallet-related HTTP requests
type WalletHandler struct {
	walletService *services.WalletService
}

// NewWalletHandler creates a new wallet handler
func NewWalletHandler(walletService *services.WalletService) *WalletHandler {
	return &WalletHandler{
		walletService: walletService,
	}
}

// ============================================
// WALLET ENDPOINTS
// ============================================

// GetWalletBalance handles GET /api/v1/wallet/balance
func (h *WalletHandler) GetWalletBalance(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code": 401,
			"msg":  "Unauthorized",
		})
	}

	balance, err := h.walletService.GetWalletBalance(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code": 500,
			"msg":  "Failed to get wallet balance",
		})
	}

	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  "success",
		"data": balance,
	})
}

// ============================================
// COIN PURCHASE ENDPOINTS
// ============================================

// GetCoinPackages handles GET /api/v1/wallet/coin-packages
func (h *WalletHandler) GetCoinPackages(c *fiber.Ctx) error {
	packages, err := h.walletService.GetCoinPackages(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code": 500,
			"msg":  "Failed to get coin packages",
		})
	}

	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  "success",
		"data": packages,
	})
}

// PurchaseCoins handles POST /api/v1/wallet/purchase-coins
func (h *WalletHandler) PurchaseCoins(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code": 401,
			"msg":  "Unauthorized",
		})
	}

	var req models.PurchaseCoinsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "Invalid request body",
		})
	}

	// Validate request
	if req.Coins <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "Coins must be greater than 0",
		})
	}

	if req.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "Amount must be greater than 0",
		})
	}

	if req.PaymentMethod == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "Payment method is required",
		})
	}

	transaction, err := h.walletService.PurchaseCoins(c.Context(), userID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code": 500,
			"msg":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  "Coins purchased successfully",
		"data": transaction,
	})
}

// ============================================
// WITHDRAWAL ENDPOINTS
// ============================================

// RequestWithdrawal handles POST /api/v1/wallet/withdraw
func (h *WalletHandler) RequestWithdrawal(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code": 401,
			"msg":  "Unauthorized",
		})
	}

	var req models.WithdrawRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "Invalid request body",
		})
	}

	// Validate request
	if req.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "Amount must be greater than 0",
		})
	}

	if req.WithdrawalMethod == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "Withdrawal method is required",
		})
	}

	if req.AccountDetails == nil || len(req.AccountDetails) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "Account details are required",
		})
	}

	withdrawalReq, err := h.walletService.RequestWithdrawal(c.Context(), userID, &req)
	if err != nil {
		statusCode := fiber.StatusInternalServerError
		if err.Error() == "insufficient balance" {
			statusCode = fiber.StatusBadRequest
		}
		return c.Status(statusCode).JSON(fiber.Map{
			"code": statusCode,
			"msg":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  "Withdrawal request submitted successfully",
		"data": withdrawalReq,
	})
}

// GetWithdrawalHistory handles GET /api/v1/wallet/withdrawal-history
func (h *WalletHandler) GetWithdrawalHistory(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code": 401,
			"msg":  "Unauthorized",
		})
	}

	// Get pagination parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))

	history, err := h.walletService.GetWithdrawalHistory(c.Context(), userID, page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code": 500,
			"msg":  "Failed to get withdrawal history",
		})
	}

	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  "success",
		"data": history,
	})
}

// ============================================
// TRANSACTION ENDPOINTS
// ============================================

// GetTransactionHistory handles GET /api/v1/wallet/transactions
func (h *WalletHandler) GetTransactionHistory(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code": 401,
			"msg":  "Unauthorized",
		})
	}

	// Get pagination parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))

	transactions, err := h.walletService.GetTransactionHistory(c.Context(), userID, page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code": 500,
			"msg":  "Failed to get transaction history",
		})
	}

	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  "success",
		"data": transactions,
	})
}

// ============================================
// PAYOUT METHOD ENDPOINTS
// ============================================

// AddPayoutMethod handles POST /api/v1/wallet/payout-methods
func (h *WalletHandler) AddPayoutMethod(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code": 401,
			"msg":  "Unauthorized",
		})
	}

	var req models.AddPayoutMethodRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "Invalid request body",
		})
	}

	// Validate request
	if req.MethodType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "Method type is required",
		})
	}

	if req.AccountName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "Account name is required",
		})
	}

	if req.AccountDetails == nil || len(req.AccountDetails) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "Account details are required",
		})
	}

	method, err := h.walletService.AddPayoutMethod(c.Context(), userID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code": 500,
			"msg":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  "Payout method added successfully",
		"data": method,
	})
}

// GetPayoutMethods handles GET /api/v1/wallet/payout-methods
func (h *WalletHandler) GetPayoutMethods(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code": 401,
			"msg":  "Unauthorized",
		})
	}

	methods, err := h.walletService.GetPayoutMethods(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code": 500,
			"msg":  "Failed to get payout methods",
		})
	}

	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  "success",
		"data": methods,
	})
}

// DeletePayoutMethod handles DELETE /api/v1/wallet/payout-methods/:id
func (h *WalletHandler) DeletePayoutMethod(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code": 401,
			"msg":  "Unauthorized",
		})
	}

	methodID := c.Params("id")
	if methodID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "Invalid payout method ID",
		})
	}

	err := h.walletService.DeletePayoutMethod(c.Context(), methodID, userID)
	if err != nil {
		statusCode := fiber.StatusInternalServerError
		if err.Error() == "payout method not found" {
			statusCode = fiber.StatusNotFound
		}
		return c.Status(statusCode).JSON(fiber.Map{
			"code": statusCode,
			"msg":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  "Payout method deleted successfully",
	})
}

// ============================================
// LEGACY ANDROID ENDPOINTS (Compatibility)
// ============================================

// ShowPayout handles POST /api/v1/showPayout (legacy)
func (h *WalletHandler) ShowPayout(c *fiber.Ctx) error {
	return h.GetPayoutMethods(c)
}

// AddPayout handles POST /api/v1/addPayout (legacy)
func (h *WalletHandler) AddPayout(c *fiber.Ctx) error {
	return h.AddPayoutMethod(c)
}

// PurchaseCoin handles POST /api/v1/purchaseCoin (legacy)
func (h *WalletHandler) PurchaseCoin(c *fiber.Ctx) error {
	return h.PurchaseCoins(c)
}

// WithdrawRequest handles POST /api/v1/withdrawRequest (legacy)
func (h *WalletHandler) WithdrawRequest(c *fiber.Ctx) error {
	return h.RequestWithdrawal(c)
}

// ShowWithdrawalHistory handles POST /api/v1/showWithdrawalHistory (legacy)
func (h *WalletHandler) ShowWithdrawalHistory(c *fiber.Ctx) error {
	return h.GetWithdrawalHistory(c)
}

// ShowOrderHistory handles POST /api/v1/showOrderHistory
func (h *WalletHandler) ShowOrderHistory(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// Get transaction history (purchases)
	history, err := h.walletService.GetTransactionHistory(c.Context(), userID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code": 500,
			"msg":  "Failed to get order history",
		})
	}

	hasMore := len(history) == limit

	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  "success",
		"data": fiber.Map{
			"orders":   history,
			"page":     page,
			"limit":    limit,
			"has_more": hasMore,
		},
	})
}
