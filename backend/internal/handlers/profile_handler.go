package handlers

import (
	"strconv"

	"lomi-backend/internal/models"
	"lomi-backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type ProfileHandler struct {
	profileService *services.ProfileService
}

func NewProfileHandler(profileService *services.ProfileService) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
	}
}

// ============================================
// PROFILE MANAGEMENT
// ============================================

// EditProfile handles POST /api/v1/editProfile
func (h *ProfileHandler) EditProfile(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	var req models.EditProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "Invalid request body",
		})
	}

	err := h.profileService.UpdateProfile(c.Context(), userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  "Profile updated successfully",
	})
}

// ============================================
// FOLLOW MANAGEMENT
// ============================================

// FollowUser handles POST /api/v1/followUser
func (h *ProfileHandler) FollowUser(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	var req models.FollowUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "Invalid request body",
		})
	}

	if req.UserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "user_id is required",
		})
	}

	if req.Action != "follow" && req.Action != "unfollow" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "action must be 'follow' or 'unfollow'",
		})
	}

	err := h.profileService.FollowUser(c.Context(), userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  err.Error(),
		})
	}

	message := "User followed successfully"
	if req.Action == "unfollow" {
		message = "User unfollowed successfully"
	}

	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  message,
	})
}

// ShowFollowers handles POST /api/v1/showFollowers
func (h *ProfileHandler) ShowFollowers(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	viewerID := claims["user_id"].(string)

	// Get user_id from request body
	var reqBody struct {
		UserID string `json:"user_id"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "Invalid request body",
		})
	}

	userID := reqBody.UserID
	if userID == "" {
		userID = viewerID // If no user_id provided, show own followers
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))

	followers, err := h.profileService.GetFollowers(c.Context(), userID, viewerID, page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code": 500,
			"msg":  "Failed to get followers",
		})
	}

	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  "success",
		"data": followers,
	})
}

// ShowFollowing handles POST /api/v1/showFollowing
func (h *ProfileHandler) ShowFollowing(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	viewerID := claims["user_id"].(string)

	// Get user_id from request body
	var reqBody struct {
		UserID string `json:"user_id"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "Invalid request body",
		})
	}

	userID := reqBody.UserID
	if userID == "" {
		userID = viewerID // If no user_id provided, show own following
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))

	following, err := h.profileService.GetFollowing(c.Context(), userID, viewerID, page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code": 500,
			"msg":  "Failed to get following",
		})
	}

	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  "success",
		"data": following,
	})
}

// ============================================
// BLOCK MANAGEMENT
// ============================================

// BlockUser handles POST /api/v1/blockUser
func (h *ProfileHandler) BlockUser(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	var req models.BlockUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "Invalid request body",
		})
	}

	if req.UserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "user_id is required",
		})
	}

	if req.Action != "block" && req.Action != "unblock" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "action must be 'block' or 'unblock'",
		})
	}

	err := h.profileService.BlockUser(c.Context(), userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  err.Error(),
		})
	}

	message := "User blocked successfully"
	if req.Action == "unblock" {
		message = "User unblocked successfully"
	}

	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  message,
	})
}

// ShowBlockedUsers handles POST /api/v1/showBlockedUsers
func (h *ProfileHandler) ShowBlockedUsers(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))

	blocked, err := h.profileService.GetBlockedUsers(c.Context(), userID, page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code": 500,
			"msg":  "Failed to get blocked users",
		})
	}

	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  "success",
		"data": blocked,
	})
}

// ============================================
// PRIVACY SETTINGS
// ============================================

// AddPrivacySetting handles POST /api/v1/addPrivacySetting
func (h *ProfileHandler) AddPrivacySetting(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	var req models.UpdatePrivacySettingsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "Invalid request body",
		})
	}

	settings, err := h.profileService.UpdatePrivacySettings(c.Context(), userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  "Privacy settings updated successfully",
		"data": settings,
	})
}

// GetPrivacySettings handles GET /api/v1/users/privacy (modern endpoint)
func (h *ProfileHandler) GetPrivacySettings(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	settings, err := h.profileService.GetPrivacySettings(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code": 500,
			"msg":  "Failed to get privacy settings",
		})
	}

	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  "success",
		"data": settings,
	})
}

// ============================================
// NOTIFICATION SETTINGS
// ============================================

// UpdatePushNotificationSettings handles POST /api/v1/updatePushNotificationSettings
func (h *ProfileHandler) UpdatePushNotificationSettings(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	var req models.UpdateNotificationSettingsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "Invalid request body",
		})
	}

	settings, err := h.profileService.UpdateNotificationSettings(c.Context(), userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  "Notification settings updated successfully",
		"data": settings,
	})
}

// GetPushNotifications handles GET /api/v1/users/push-notifications (modern endpoint)
func (h *ProfileHandler) GetPushNotifications(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	settings, err := h.profileService.GetNotificationSettings(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code": 500,
			"msg":  "Failed to get notification settings",
		})
	}

	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  "success",
		"data": settings,
	})
}

// ============================================
// REFERRAL SYSTEM
// ============================================

// GetReferralCode handles GET /api/v1/getReferralCode
func (h *ProfileHandler) GetReferralCode(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	stats, err := h.profileService.GetReferralCode(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code": 500,
			"msg":  "Failed to get referral code",
		})
	}

	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  "success",
		"data": stats,
	})
}

// ApplyReferralCode handles POST /api/v1/applyReferralCode
func (h *ProfileHandler) ApplyReferralCode(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	var req models.ApplyReferralCodeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "Invalid request body",
		})
	}

	if req.ReferralCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "referral_code is required",
		})
	}

	err := h.profileService.ApplyReferralCode(c.Context(), userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  "Referral code applied successfully. You and your referrer have been rewarded!",
	})
}

// ============================================
// ACCOUNT MANAGEMENT
// ============================================

// DeleteUserAccount handles POST /api/v1/deleteUserAccount
func (h *ProfileHandler) DeleteUserAccount(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	err := h.profileService.DeleteUserAccount(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  "Account deleted successfully",
	})
}

// UserVerificationRequest handles POST /api/v1/userVerificationRequest
func (h *ProfileHandler) UserVerificationRequest(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	var req models.VerificationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "Invalid request body",
		})
	}

	if req.SelfieURL == "" || req.IDDocumentURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "selfie_url and id_document_url are required",
		})
	}

	err := h.profileService.RequestVerification(c.Context(), userID, req.SelfieURL, req.IDDocumentURL)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  "Verification request submitted successfully. We'll review it within 24-48 hours.",
	})
}

// ReportUser handles POST /api/v1/reportUser
func (h *ProfileHandler) ReportUser(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	reporterID := claims["user_id"].(string)

	var req models.ReportUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "Invalid request body",
		})
	}

	if req.UserID == "" || req.Reason == "" || req.Description == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "user_id, reason, and description are required",
		})
	}

	err := h.profileService.ReportUser(c.Context(), reporterID, req.UserID, req.Reason, req.Description, req.Screenshots)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  "Report submitted successfully. We'll review it and take appropriate action.",
	})
}

// ============================================
// SOCIAL FEATURES
// ============================================

// GenerateQRCode handles GET /api/v1/generateQRCode
func (h *ProfileHandler) GenerateQRCode(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	qrData, err := h.profileService.GenerateQRCode(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code": 500,
			"msg":  "Failed to generate QR code",
		})
	}

	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  "success",
		"data": qrData,
	})
}

// ShareProfile handles POST /api/v1/shareProfile
func (h *ProfileHandler) ShareProfile(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	sharedBy := claims["user_id"].(string)

	var req models.ShareProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "Invalid request body",
		})
	}

	if req.UserID == "" || req.Platform == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "user_id and platform are required",
		})
	}

	err := h.profileService.ShareProfile(c.Context(), sharedBy, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  "Profile share tracked successfully",
	})
}

// ============================================
// APP SETTINGS
// ============================================

// ChangeAppLanguage handles POST /api/v1/changeAppLanguage
func (h *ProfileHandler) ChangeAppLanguage(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)
	
	var req models.ChangeLanguageRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "Invalid request body",
		})
	}
	
	if req.Language == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "language is required",
		})
	}
	
	err := h.profileService.ChangeAppLanguage(c.Context(), userID, req.Language)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  err.Error(),
		})
	}
	
	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  "Language changed successfully",
	})
}

// ChangeAppTheme handles POST /api/v1/changeAppTheme
func (h *ProfileHandler) ChangeAppTheme(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)
	
	var req models.ChangeThemeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "Invalid request body",
		})
	}
	
	if req.Theme == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "theme is required",
		})
	}
	
	err := h.profileService.ChangeAppTheme(c.Context(), userID, req.Theme)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  err.Error(),
		})
	}
	
	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  "Theme changed successfully",
	})
}

// ClearCache handles POST /api/v1/clearCache
func (h *ProfileHandler) ClearCache(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)
	
	err := h.profileService.ClearCache(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code": 500,
			"msg":  "Failed to clear cache",
		})
	}
	
	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  "Cache cleared successfully",
	})
}
