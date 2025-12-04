package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// LegacyHandler handles requests from the legacy Android app
type LegacyHandler struct{}

func NewLegacyHandler() *LegacyHandler {
	return &LegacyHandler{}
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
		"msg":  "success",
		"msg_array": fiber.Map{
			"count": "0",
		},
	})
}

// CheckPhoneNo handles /api/checkPhoneNo
func (h *LegacyHandler) CheckPhoneNo(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  "success",
	})
}
