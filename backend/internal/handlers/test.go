package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// TestEndpoint is a simple endpoint to verify backend is reachable
func TestEndpoint(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"message": "Backend is reachable! ✅",
		"method":  c.Method(),
		"path":    c.Path(),
		"ip":      c.IP(),
		"headers": c.GetReqHeaders(),
	})
}

// TestAuthEndpoint tests if auth endpoint is reachable
func TestAuthEndpoint(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"message": "Auth endpoint is reachable! ✅",
		"method":  c.Method(),
		"path":    c.Path(),
		"note":    "This endpoint accepts POST requests with Authorization: tma <initData>",
	})
}

