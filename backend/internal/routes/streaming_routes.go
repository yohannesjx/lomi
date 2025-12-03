package routes

import (
	"lomi-backend/config"
	"lomi-backend/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

// SetupStreamingRoutes adds TikTok-style video/streaming endpoints
// These routes match the exact API contract from the Android/iOS apps
func SetupStreamingRoutes(app *fiber.App) {
	// Create streaming handler
	streamingHandler := handlers.NewStreamingHandler(config.Cfg)

	// TikTok API base path: /api/
	api := app.Group("/api")

	// ==================== PUBLIC ENDPOINTS ====================
	// These endpoints don't require authentication (for initial login)

	// 1. POST /api/checkUsername - Check username availability
	api.Post("/checkUsername", streamingHandler.CheckUsername)

	// 2. POST /api/registerUser - Social login/signup
	api.Post("/registerUser", streamingHandler.RegisterUser)

	// ==================== PROTECTED ENDPOINTS ====================
	// These endpoints require authentication
	// Note: The TikTok apps send auth_token in the request body, not headers
	// So we'll handle auth validation inside each handler for now

	// 2. POST /api/showUserDetail - Get user profile & wallet
	api.Post("/showUserDetail", streamingHandler.ShowUserDetail)

	// 3. POST /api/showRelatedVideos - Home feed (For You page)
	api.Post("/showRelatedVideos", streamingHandler.ShowRelatedVideos)

	// 4. POST /api/liveStream - Start live streaming
	api.Post("/liveStream", streamingHandler.LiveStream)

	// 5. POST /api/sendGift - Send virtual gift
	api.Post("/sendGift", streamingHandler.SendGift)

	// 6. POST /api/purchaseCoin - Buy coins
	api.Post("/purchaseCoin", streamingHandler.PurchaseCoin)

	// ==================== BONUS ENDPOINTS ====================
	// Additional endpoints that might be needed

	// Show coin packages
	api.Post("/showCoinWorth", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"code": 200,
			"msg": fiber.Map{
				"CoinPackages": []fiber.Map{
					{
						"id":    1,
						"coins": 100,
						"price": 0.99,
						"title": "Starter Pack",
					},
					{
						"id":    2,
						"coins": 500,
						"price": 4.99,
						"title": "Popular Pack",
					},
					{
						"id":    3,
						"coins": 1000,
						"price": 9.99,
						"title": "Best Value",
					},
					{
						"id":    4,
						"coins": 2500,
						"price": 19.99,
						"title": "Premium Pack",
					},
					{
						"id":    5,
						"coins": 5000,
						"price": 39.99,
						"title": "Ultimate Pack",
					},
				},
			},
		})
	})

	// Show gifts catalog
	api.Post("/showGifts", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"code": 200,
			"msg": []fiber.Map{
				{
					"Gift": fiber.Map{
						"id":    1,
						"title": "Rose",
						"image": "https://cdn.example.com/gifts/rose.png",
						"coin":  10,
						"type":  "normal",
					},
				},
				{
					"Gift": fiber.Map{
						"id":    2,
						"title": "Heart",
						"image": "https://cdn.example.com/gifts/heart.png",
						"coin":  50,
						"type":  "normal",
					},
				},
				{
					"Gift": fiber.Map{
						"id":    3,
						"title": "Diamond",
						"image": "https://cdn.example.com/gifts/diamond.png",
						"coin":  1000,
						"type":  "premium",
					},
				},
				{
					"Gift": fiber.Map{
						"id":    4,
						"title": "Universe",
						"image": "https://cdn.example.com/gifts/universe.png",
						"coin":  5000,
						"type":  "luxury",
					},
				},
			},
		})
	})
}
