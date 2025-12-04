package handlers

import (
	"lomi-backend/config"
	"lomi-backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// LegacyHandler handles requests from the legacy Android app
type LegacyHandler struct {
	profileService *services.ProfileService
}

func NewLegacyHandler(profileService *services.ProfileService) *LegacyHandler {
	return &LegacyHandler{
		profileService: profileService,
	}
}

// ShowRooms handles /api/showRooms
func (h *LegacyHandler) ShowRooms(c *fiber.Ctx) error {
	// Return empty list of rooms for now
	return c.JSON(fiber.Map{
		"code":      200,
		"msg":       "success",
		"msg_array": []interface{}{},
	})
}

// ShowFriendsStories handles /api/showFriendsStories
func (h *LegacyHandler) ShowFriendsStories(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"code":      200,
		"msg":       "success",
		"msg_array": []interface{}{},
	})
}

// ShowSettings handles /api/showSettings
func (h *LegacyHandler) ShowSettings(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  "success",
		"msg_array": fiber.Map{
			"privacy":           "1",
			"push_notification": "1",
		},
	})
}

// ShowVideoDetailAd handles /api/showVideoDetailAd
func (h *LegacyHandler) ShowVideoDetailAd(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"code":      200,
		"msg":       "success",
		"msg_array": []interface{}{},
	})
}

// ShowUnReadNotifications handles /api/showUnReadNotifications
func (h *LegacyHandler) ShowUnReadNotifications(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  "0", // Return count as string in msg field
	})
}

// CheckPhoneNo handles /api/checkPhoneNo
func (h *LegacyHandler) CheckPhoneNo(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  "success",
	})
}

// ShowUserDetail handles /api/showUserDetail
func (h *LegacyHandler) ShowUserDetail(c *fiber.Ctx) error {
	var req struct {
		UserID    string `json:"user_id"`
		AuthToken string `json:"auth_token"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "Invalid request body",
		})
	}

	if req.UserID == "" && req.AuthToken != "" {
		// Parse token
		token, err := jwt.Parse(req.AuthToken, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Cfg.JWTSecret), nil
		})

		if err == nil && token.Valid {
			claims, ok := token.Claims.(jwt.MapClaims)
			if ok {
				if id, ok := claims["user_id"].(string); ok {
					req.UserID = id
				}
			}
		}
	}

	if req.UserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "user_id is required",
		})
	}

	userDetail, err := h.profileService.GetUserDetail(c.Context(), req.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code": 500,
			"msg":  "Failed to get user details",
		})
	}

	response := fiber.Map{
		"User": userDetail,
	}

	if userDetail.PrivacySetting != nil {
		response["PrivacySetting"] = userDetail.PrivacySetting
	}
	if userDetail.PushSetting != nil {
		response["PushNotification"] = userDetail.PushSetting
	}

	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  response,
	})
}
