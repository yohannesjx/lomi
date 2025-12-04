package routes

import (
	"lomi-backend/config"
	"lomi-backend/internal/handlers"
	"lomi-backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func SetupRoutes(app *fiber.App, walletHandler *handlers.WalletHandler, profileHandler *handlers.ProfileHandler, videoHandler *handlers.VideoHandler) {
	api := app.Group("/api/v1")

	// Health Check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "message": "Lomi Backend is running 🍋"})
	})

	// Legacy endpoints (now under v1)
	legacyHandler := handlers.NewLegacyHandler()
	api.Post("/showRooms", legacyHandler.ShowRooms)
	api.Post("/showFriendsStories", legacyHandler.ShowFriendsStories)
	api.Post("/showSettings", legacyHandler.ShowSettings)
	api.Post("/showVideoDetailAd", legacyHandler.ShowVideoDetailAd)
	api.Post("/showUnReadNotifications", legacyHandler.ShowUnReadNotifications)
	api.Post("/checkPhoneNo", legacyHandler.CheckPhoneNo)
	api.Post("/showUserDetail", legacyHandler.ShowUserDetail)

	// Test endpoints for debugging
	api.Get("/test", handlers.TestEndpoint)
	api.Post("/test", handlers.TestEndpoint)
	api.Get("/test/auth", handlers.TestAuthEndpoint)
	api.Post("/test/auth", handlers.TestAuthEndpoint)
	api.Get("/test/s3", handlers.TestS3Connection)
	api.Get("/test/jwt", handlers.TestGetJWT)  // Generate JWT for user by ID (testing only)
	api.Post("/test/jwt", handlers.TestGetJWT) // Generate JWT for user by ID (testing only)

	// Public routes
	authHandler := handlers.NewAuthHandler(config.Cfg)

	// Telegram Mini App login (initData method) - Auto-authenticates on app open
	api.Post("/auth/telegram", authHandler.TelegramLogin)
	api.Post("/auth/google", authHandler.GoogleLogin)

	api.Post("/auth/refresh", handlers.RefreshToken) // TODO: Implement

	// ============================================
	// PUBLIC WALLET ENDPOINTS (No Auth Required)
	// ============================================

	// Coin Packages (Public - anyone can view packages)
	api.Get("/wallet/coin-packages", walletHandler.GetCoinPackages)

	// Protected routes (require authentication)
	protected := api.Group("", middleware.AuthMiddleware)

	// Test media upload (protected)
	protected.Get("/test/media-upload", handlers.TestMediaUpload)

	// User Profile
	protected.Get("/users/me", handlers.GetMe)
	protected.Put("/users/me", handlers.UpdateProfile)
	protected.Get("/users", handlers.GetAllUsers) // List all users (for testing)

	// Onboarding
	protected.Get("/onboarding/status", handlers.GetOnboardingStatus)
	protected.Patch("/onboarding/progress", handlers.UpdateOnboardingProgress)

	// ============================================
	// PROFILE MANAGEMENT (Phase 1)
	// ============================================

	// Profile Updates
	protected.Post("/editProfile", profileHandler.EditProfile)

	// Follow System
	protected.Post("/followUser", profileHandler.FollowUser)
	protected.Post("/showFollowers", profileHandler.ShowFollowers)
	protected.Post("/showFollowing", profileHandler.ShowFollowing)

	// Block System
	protected.Post("/blockUser", profileHandler.BlockUser)
	protected.Post("/showBlockedUsers", profileHandler.ShowBlockedUsers)

	// Privacy Settings (Legacy + Modern)
	protected.Post("/addPrivacySetting", profileHandler.AddPrivacySetting)
	protected.Get("/users/privacy", profileHandler.GetPrivacySettings)
	protected.Post("/users/privacy", profileHandler.AddPrivacySetting)

	// Notification Settings (Legacy + Modern)
	protected.Post("/updatePushNotificationSettings", profileHandler.UpdatePushNotificationSettings)
	protected.Get("/users/push-notifications", profileHandler.GetPushNotifications)
	protected.Post("/users/push-notifications", profileHandler.UpdatePushNotificationSettings)

	// Referral System
	protected.Get("/getReferralCode", profileHandler.GetReferralCode)
	protected.Post("/applyReferralCode", profileHandler.ApplyReferralCode)

	// ============================================
	// SOCIAL CONTENT (Phase 1.2)
	// ============================================

	// User Videos
	protected.Post("/showVideosAgainstUserID", videoHandler.ShowVideosAgainstUserID)
	protected.Post("/showUserLikedVideos", videoHandler.ShowUserLikedVideos)
	protected.Post("/showUserRepostedVideos", videoHandler.ShowUserRepostedVideos)
	protected.Post("/showFavouriteVideos", videoHandler.ShowFavouriteVideos)

	// ============================================
	// ACCOUNT MANAGEMENT (Phase 2.3)
	// ============================================

	protected.Post("/deleteUserAccount", profileHandler.DeleteUserAccount)
	protected.Post("/userVerificationRequest", profileHandler.UserVerificationRequest)
	protected.Post("/reportUser", profileHandler.ReportUser)

	// ============================================
	// SOCIAL FEATURES (Phase 4)
	// ============================================

	protected.Get("/generateQRCode", profileHandler.GenerateQRCode)
	protected.Post("/shareProfile", profileHandler.ShareProfile)

	// Settings (keeping existing handlers for now)

	// Media
	protected.Post("/users/media", handlers.UploadMedia)
	protected.Post("/users/media/upload-complete", handlers.UploadComplete) // Batch upload completion
	protected.Get("/users/media/moderation-status", handlers.GetMyModerationStatus)
	protected.Get("/users/:user_id/media", handlers.GetUserMedia)
	protected.Delete("/users/media/:id", handlers.DeleteMedia)
	protected.Get("/users/media/upload-url", handlers.GetPresignedUploadURL)

	// Discovery & Swiping (with rate limiting)
	protected.Get("/discover/swipe", handlers.GetSwipeCards)
	protected.Post("/discover/swipe", middleware.SwipeRateLimit(), handlers.SwipeAction)
	protected.Get("/discover/feed", handlers.GetExploreFeed)

	// Matches
	protected.Get("/matches", handlers.GetMatches)
	protected.Get("/matches/:id", handlers.GetMatchDetails)
	protected.Delete("/matches/:id", handlers.Unmatch)

	// Chat (with rate limiting for messages)
	protected.Get("/chats", handlers.GetChats)
	protected.Get("/chats/:id/messages", handlers.GetMessages)
	protected.Post("/chats/:id/messages", middleware.MessageRateLimit(), handlers.SendMessage)
	protected.Put("/chats/:id/read", handlers.MarkMessagesAsRead)

	// Gifts (Luxury System)
	protected.Get("/gifts/shop", handlers.GetGiftShop)
	protected.Post("/gifts/send", handlers.SendGiftLuxury)
	protected.Get("/gifts/received", handlers.GetGiftsReceived)

	// Legacy gifts endpoint (keep for backward compatibility)
	protected.Get("/gifts", handlers.GetGifts)

	// Wallet (Luxury System)
	protected.Get("/wallet/balance", handlers.GetWalletBalance)
	protected.Post("/wallet/buy", middleware.PurchaseRateLimit(), handlers.BuyCoins)
	protected.Post("/wallet/buy/webhook", handlers.CoinPurchaseWebhook) // Telebirr webhook

	// Legacy coins endpoints (keep for backward compatibility)
	protected.Get("/coins/balance", handlers.GetCoinBalance)
	protected.Post("/coins/purchase", middleware.PurchaseRateLimit(), handlers.PurchaseCoins)
	protected.Post("/coins/purchase/confirm", handlers.ConfirmCoinPurchase) // Webhook endpoint
	protected.Get("/coins/transactions", handlers.GetCoinTransactions)

	// Cashout (Luxury System)
	protected.Post("/cashout/request", handlers.RequestCashout)

	// Legacy payouts endpoints (keep for backward compatibility)
	protected.Get("/payouts/balance", handlers.GetPayoutBalance)
	protected.Post("/payouts/request", handlers.RequestPayout)
	protected.Get("/payouts/history", handlers.GetPayoutHistory)

	// ============================================
	// NEW WALLET MANAGEMENT SYSTEM (Production-Grade)
	// ============================================

	// Wallet Balance & Info
	protected.Get("/wallet/v2/balance", walletHandler.GetWalletBalance)

	// Coin Purchase
	protected.Post("/wallet/purchase-coins", walletHandler.PurchaseCoins)

	// Withdrawals
	protected.Post("/wallet/withdraw", walletHandler.RequestWithdrawal)
	protected.Get("/wallet/withdrawal-history", walletHandler.GetWithdrawalHistory)

	// Transactions
	protected.Get("/wallet/transactions", walletHandler.GetTransactionHistory)

	// Payout Methods
	protected.Post("/wallet/payout-methods", walletHandler.AddPayoutMethod)
	protected.Get("/wallet/payout-methods", walletHandler.GetPayoutMethods)
	protected.Delete("/wallet/payout-methods/:id", walletHandler.DeletePayoutMethod)

	// Legacy Android Endpoints (Backward Compatibility)
	protected.Post("/showPayout", walletHandler.ShowPayout)
	protected.Post("/addPayout", walletHandler.AddPayout)
	protected.Post("/purchaseCoin", walletHandler.PurchaseCoin)
	protected.Post("/withdrawRequest", walletHandler.WithdrawRequest)
	protected.Post("/showWithdrawalHistory", walletHandler.ShowWithdrawalHistory)

	// Verification
	protected.Post("/verification/submit", handlers.SubmitVerification)
	protected.Get("/verification/status", handlers.GetVerificationStatus)

	// Reports & Blocks
	protected.Post("/reports", handlers.ReportUser)
	protected.Post("/reports/photo", handlers.ReportPhoto)
	protected.Post("/blocks", handlers.BlockUser)
	protected.Delete("/blocks/:user_id", handlers.UnblockUser)
	protected.Get("/blocks", handlers.GetBlockedUsers)

	// Reward Channels
	protected.Get("/coins/earn/channels", handlers.GetRewardChannels)
	protected.Post("/coins/earn/claim", handlers.ClaimChannelReward)

	// Leaderboard
	protected.Get("/leaderboard/top-gifted", handlers.GetTopGiftedUsers)

	// Who Likes You (Likes Reveal)
	protected.Get("/likes/pending", handlers.GetPendingLikes)
	protected.Post("/likes/reveal", handlers.RevealLike)

	// Admin routes (should have admin middleware in production)
	admin := protected.Group("/admin")
	admin.Get("/reports/pending", handlers.GetPendingReports)
	admin.Put("/reports/:id/review", handlers.ReviewReport)
	admin.Get("/payouts/pending", handlers.GetPendingPayouts)
	admin.Put("/payouts/:id/process", handlers.ProcessPayout)

	// Photo Moderation Monitoring (Phase 3)
	admin.Get("/queue-stats", handlers.GetQueueStats)
	admin.Get("/moderation/dashboard", handlers.GetModerationDashboard)
	admin.Put("/moderation/rejected/:id/verify", handlers.VerifyRejectedPhoto)
	admin.Delete("/moderation/rejected/:id", handlers.DeleteRejectedPhoto)

	// WebSocket - Legacy (keep for backward compatibility)
	api.Get("/ws", websocket.New(handlers.HandleWebSocket))

	// WebSocket - Unified Chat (handles both private 1-on-1 and live streaming chat)
	api.Get("/ws/chat", websocket.New(handlers.HandleUnifiedChat))

	// Live Chat HTTP Endpoints
	protected.Get("/live/:id/viewers", handlers.GetLiveViewerCount)
	protected.Get("/live/:id/pinned", handlers.GetPinnedMessage)
}
