package handlers

import (
	"lomi-backend/internal/database"
	"lomi-backend/internal/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// GetPendingReports returns all pending reports for admin review
func GetPendingReports(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 50)
	offset := (page - 1) * limit

	var reports []models.Report
	if err := database.DB.Where("is_reviewed = ?", false).
		Preload("Reporter").
		Preload("ReportedUser").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&reports).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch reports"})
	}

	return c.JSON(fiber.Map{
		"reports": reports,
		"page":    page,
		"limit":   limit,
	})
}

// ReviewReport allows admin to review and take action on a report
func ReviewReport(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	adminIDStr := claims["user_id"].(string)
	adminID, _ := uuid.Parse(adminIDStr)

	reportID := c.Params("id")
	var req struct {
		Action      string `json:"action"` // "approve", "reject", "warn", "ban"
		ActionTaken string `json:"action_taken,omitempty"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	var report models.Report
	if err := database.DB.First(&report, "id = ?", reportID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Report not found"})
	}

	now := time.Now()
	report.IsReviewed = true
	report.ReviewedBy = &adminID
	report.ReviewedAt = &now
	report.ActionTaken = req.ActionTaken

	// Take action based on review
	switch req.Action {
	case "ban":
		// Ban the reported user
		database.DB.Model(&models.User{}).
			Where("id = ?", report.ReportedUserID).
			Update("is_active", false)
	case "warn":
		// Add warning to user (could be stored in user metadata)
		// For now, just mark as reviewed
	}

	if err := database.DB.Save(&report).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update report"})
	}

	return c.JSON(fiber.Map{
		"message": "Report reviewed successfully",
		"report":  report,
	})
}

// GetPendingPayouts returns all pending payout requests for admin review
func GetPendingPayouts(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 50)
	offset := (page - 1) * limit

	var payouts []models.Payout
	if err := database.DB.Where("status = ?", models.PayoutStatusPending).
		Preload("User").
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

// ProcessPayout allows admin to approve or reject a payout request
func ProcessPayout(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	adminIDStr := claims["user_id"].(string)
	adminID, _ := uuid.Parse(adminIDStr)

	payoutID := c.Params("id")
	var req struct {
		Action            string `json:"action"` // "approve", "reject"
		PaymentReference  string `json:"payment_reference,omitempty"`
		RejectionReason   string `json:"rejection_reason,omitempty"`
		AdminNotes        string `json:"admin_notes,omitempty"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	var payout models.Payout
	if err := database.DB.First(&payout, "id = ?", payoutID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Payout not found"})
	}

	if payout.Status != models.PayoutStatusPending {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Payout already processed"})
	}

	now := time.Now()
	payout.ProcessedBy = &adminID
	payout.ProcessedAt = &now
	payout.AdminNotes = req.AdminNotes

	if req.Action == "approve" {
		payout.Status = models.PayoutStatusProcessing
		payout.PaymentReference = req.PaymentReference

		// TODO: Integrate with payment gateway to process payout
		// For now, mark as processing. In production, you'd:
		// 1. Call Telebirr/CBE Birr API to send money
		// 2. Update status to "completed" on success
		// 3. Update status to "rejected" on failure

	} else if req.Action == "reject" {
		payout.Status = models.PayoutStatusRejected
		payout.RejectionReason = req.RejectionReason

		// Refund the amount back to user's gift balance
		var user models.User
		if err := database.DB.First(&user, "id = ?", payout.UserID).Error; err == nil {
			user.GiftBalance += payout.GiftBalanceAmount
			database.DB.Save(&user)
		}
	}

	if err := database.DB.Save(&payout).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update payout"})
	}

	return c.JSON(fiber.Map{
		"message": "Payout processed successfully",
		"payout":  payout,
	})
}

