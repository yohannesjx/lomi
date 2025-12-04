package main

import (
	"log"
	"lomi-backend/config"
	"lomi-backend/internal/database"
	"lomi-backend/internal/handlers"
	"lomi-backend/internal/repositories"
	"lomi-backend/internal/routes"
	"lomi-backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// 1. Load Configuration
	cfg := config.LoadConfig()

	// 2. Connect to Database (GORM)
	database.ConnectDB(cfg)

	// 3. Connect to Database (sqlx for wallet system)
	database.ConnectSqlxDB(cfg)

	// 4. Connect to Redis
	database.ConnectRedis(cfg)

	// 5. Connect to S3/R2
	database.ConnectS3(cfg)

	// 6. Initialize Notification Service
	services.InitNotificationService(
		cfg.TelegramBotToken,
		cfg.OneSignalAppID,
		cfg.OneSignalAPIKey,
		cfg.FirebaseServerKey,
	)

	// 7. Initialize Wallet Dependencies
	walletRepo := repositories.NewWalletRepository(database.SqlxDB)
	walletService := services.NewWalletService(walletRepo)
	walletHandler := handlers.NewWalletHandler(walletService)

	// 8. Initialize Profile Dependencies
	profileRepo := repositories.NewProfileRepository(database.SqlxDB)
	profileService := services.NewProfileService(profileRepo)
	profileHandler := handlers.NewProfileHandler(profileService)

	// 9. Initialize Video Dependencies
	videoRepo := repositories.NewVideoRepository(database.SqlxDB)
	videoService := services.NewVideoService(videoRepo)
	videoHandler := handlers.NewVideoHandler(videoService)

	// 10. Initialize Fiber App
	app := fiber.New(fiber.Config{
		AppName:      cfg.AppName,
		ServerHeader: "Lomi-Social",
		Prefork:      false, // Set to true for production if needed
	})

	// 9. Middleware
	app.Use(logger.New())  // Request logging
	app.Use(recover.New()) // Panic recovery
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Allow all for dev, restrict in prod
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, HEAD, PUT, DELETE, PATCH",
	}))

	// 11. Routes
	routes.SetupRoutes(app, walletHandler, profileHandler, videoHandler)
	routes.SetupStreamingRoutes(app) // TikTok-style streaming endpoints

	// 11. Start Server
	log.Printf("🚀 Server starting on port %s", cfg.AppPort)
	if err := app.Listen(":" + cfg.AppPort); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}

func setupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	// Health Check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Lomi Social API is running 💚",
		})
	})
}
