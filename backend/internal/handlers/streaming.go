package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"lomi-backend/config"
	"lomi-backend/internal/database"
	"lomi-backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StreamingHandler handles TikTok-style video and live streaming endpoints
type StreamingHandler struct {
	cfg *config.Config
}

func NewStreamingHandler(cfg *config.Config) *StreamingHandler {
	return &StreamingHandler{cfg: cfg}
}

// ==================== TikTok API Response Format ====================
// All responses must match the exact format: {"code": 200, "msg": {...}}

type TikTokResponse struct {
	Code int         `json:"code"`
	Msg  interface{} `json:"msg"`
}

func tikTokSuccess(msg interface{}) TikTokResponse {
	return TikTokResponse{Code: 200, Msg: msg}
}

func tikTokError(code int, message string) TikTokResponse {
	return TikTokResponse{Code: code, Msg: message}
}

// ==================== 1. POST /api/checkUsername ====================
// Check if username is available
func (h *StreamingHandler) CheckUsername(c *fiber.Ctx) error {
	var req struct {
		Username string `json:"username"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.JSON(tikTokError(400, "Invalid request body"))
	}

	if req.Username == "" {
		return c.JSON(tikTokError(400, "Username is required"))
	}

	var count int64
	if err := database.DB.Model(&models.User{}).Where("LOWER(username) = LOWER(?)", req.Username).Count(&count).Error; err != nil {
		return c.JSON(tikTokError(500, "Database error"))
	}

	if count > 0 {
		return c.JSON(tikTokError(201, "Username is already taken"))
	}

	return c.JSON(tikTokSuccess("Username is available"))
}

// ==================== 2. POST /api/registerUser ====================
// Social login endpoint - creates or returns existing user
func (h *StreamingHandler) RegisterUser(c *fiber.Ctx) error {
	var req struct {
		Username    string `json:"username"`
		DOB         string `json:"dob"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Email       string `json:"email"`
		Phone       string `json:"phone"`
		SocialID    string `json:"social_id"`
		AuthToken   string `json:"auth_token"`
		DeviceToken string `json:"device_token"`
		Social      string `json:"social"` // "google", "facebook", "apple"
	}

	if err := c.BodyParser(&req); err != nil {
		return c.JSON(tikTokError(400, "Invalid request body"))
	}

	// Validate required fields
	if req.SocialID == "" || req.Social == "" {
		return c.JSON(tikTokError(400, "social_id and social are required"))
	}

	var user models.User
	var isNewUser bool

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Try to find existing user by auth provider
		var authProvider models.AuthProvider
		result := tx.Where("provider = ? AND provider_id = ?", req.Social, req.SocialID).First(&authProvider)

		if result.Error == nil {
			// User exists, fetch full user record
			if err := tx.Where("id = ?", authProvider.UserID).First(&user).Error; err != nil {
				return err
			}
			isNewUser = false
		} else if result.Error == gorm.ErrRecordNotFound {
			// Create new user
			fullName := req.FirstName
			if req.LastName != "" {
				fullName += " " + req.LastName
			}
			if fullName == "" {
				fullName = "User"
			}

			// Generate username from email or social ID
			username := ""
			if req.Email != "" {
				username = strings.Split(req.Email, "@")[0]
			} else {
				username = "user_" + req.SocialID[:8]
			}
			username = strings.ReplaceAll(username, ".", "_")
			username = strings.ReplaceAll(username, "+", "_")

			// Make username unique
			var existingUser models.User
			baseUsername := username
			counter := 1
			for {
				if err := tx.Where("username = ?", username).First(&existingUser).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						break
					}
					return err
				}
				username = fmt.Sprintf("%s%d", baseUsername, counter)
				counter++
			}

			// Set email to nil if empty to avoid unique constraint violation
			var emailPtr *string
			if req.Email != "" {
				emailPtr = &req.Email
			}

			user = models.User{
				Username:           username,
				Name:               fullName,
				Email:              req.Email, // Can be empty
				Age:                18,
				Gender:             models.GenderOther,
				City:               "Not Set",
				RelationshipGoal:   models.GoalDating,
				Religion:           models.ReligionNone,
				VerificationStatus: models.VerificationPending,
				IsActive:           true,
				IsVerified:         false,
				Languages:          models.JSONStringArray{},
				Interests:          models.JSONStringArray{},
				Preferences:        models.JSONMap{},
				CoinBalance:        0,
				GiftBalance:        0.0,
			}

			if err := tx.Create(&user).Error; err != nil {
				return err
			}

			// Create auth provider entry
			provider := models.AuthProvider{
				UserID:     user.ID,
				Provider:   req.Social,
				ProviderID: req.SocialID,
				Email:      req.Email,
				LinkedAt:   time.Now(),
			}
			if err := tx.Create(&provider).Error; err != nil {
				return err
			}

			isNewUser = true
			log.Printf("✅ Created new user via %s: ID=%s, Username=%s, Email=%s", req.Social, user.ID, username, req.Email)
		} else {
			return result.Error
		}

		return nil
	})

	if err != nil {
		log.Printf("❌ RegisterUser error: %v", err)
		return c.JSON(tikTokError(500, "Could not process registration"))
	}

	// Generate auth token for the user
	authToken, err := generateAuthToken(user.ID)
	if err != nil {
		log.Printf("❌ Failed to generate auth token: %v", err)
		return c.JSON(tikTokError(500, "Could not generate auth token"))
	}

	// Return TikTok-style response
	response := tikTokSuccess(fiber.Map{
		"User": fiber.Map{
			"id":                   user.ID,
			"username":             req.Username,
			"first_name":           req.FirstName,
			"last_name":            req.LastName,
			"email":                user.Email,
			"phone":                req.Phone,
			"profile_pic":          "",
			"profile_pic_small":    "",
			"auth_token":           authToken,
			"wallet":               user.CoinBalance,
			"total_all_time_coins": user.TotalEarned,
			"verified":             boolToInt(user.IsVerified),
			"online":               1,
			"created":              user.CreatedAt.Format("2006-01-02 15:04:05"),
		},
		"PushNotification": fiber.Map{
			"id":              1,
			"likes":           1,
			"comments":        1,
			"new_followers":   1,
			"mentions":        1,
			"direct_messages": 1,
			"video_updates":   1,
		},
		"PrivacySetting": fiber.Map{
			"id":              1,
			"videos_download": 0,
			"direct_message":  0,
			"duet":            1,
			"liked_videos":    0,
			"video_comment":   1,
			"order_history":   0,
		},
	})

	if isNewUser {
		log.Printf("✅ New user registered: %s", user.ID)
	} else {
		log.Printf("✅ Existing user logged in: %s", user.ID)
	}

	return c.JSON(response)
}

// ==================== 2. POST /api/showUserDetail ====================
// Get user profile and wallet balance
func (h *StreamingHandler) ShowUserDetail(c *fiber.Ctx) error {
	var req struct {
		AuthToken   string `json:"auth_token"`
		UserID      string `json:"user_id"`
		OtherUserID string `json:"other_user_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.JSON(tikTokError(400, "Invalid request body"))
	}

	// Get user ID from context (set by auth middleware) or from request
	userID := c.Locals("user_id")
	if userID == nil && req.UserID != "" {
		parsedID, err := uuid.Parse(req.UserID)
		if err == nil {
			userID = parsedID
		}
	}

	if userID == nil {
		return c.JSON(tikTokError(401, "Unauthorized"))
	}

	var user models.User
	targetUserID := userID.(uuid.UUID)

	// If requesting other user's profile
	if req.OtherUserID != "" {
		otherID, err := uuid.Parse(req.OtherUserID)
		if err != nil {
			return c.JSON(tikTokError(400, "Invalid other_user_id"))
		}
		targetUserID = otherID
	}

	if err := database.DB.Where("id = ?", targetUserID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(tikTokError(404, "User not found"))
		}
		return c.JSON(tikTokError(500, "Database error"))
	}

	// Determine relationship status (following/friends/follow back)
	buttonStatus := "follow"
	if targetUserID != userID.(uuid.UUID) {
		// TODO: Check if users follow each other
		// For now, default to "follow"
		buttonStatus = "follow"
	}

	response := tikTokSuccess(fiber.Map{
		"User": fiber.Map{
			"id":                   user.ID,
			"username":             user.Name,
			"first_name":           user.Name,
			"wallet":               user.CoinBalance,
			"total_all_time_coins": user.TotalEarned,
			"verified":             boolToInt(user.IsVerified),
			"button":               buttonStatus,
			"profile_pic":          "",
			"bio":                  user.Bio,
			"city":                 user.City,
			"age":                  user.Age,
			"gender":               user.Gender,
			"online":               boolToInt(user.IsOnline),
		},
	})

	return c.JSON(response)
}

// ==================== 3. POST /api/showRelatedVideos ====================
// Home feed (For You page) - returns dummy videos for now
func (h *StreamingHandler) ShowRelatedVideos(c *fiber.Ctx) error {
	var req struct {
		UserID            string  `json:"user_id"`
		DeviceID          string  `json:"device_id"`
		StartingPoint     int     `json:"starting_point"`
		Lat               float64 `json:"lat"`
		Long              float64 `json:"long"`
		TagProduct        int     `json:"tag_product"`
		DeliveryAddressID int     `json:"delivery_address_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.JSON(tikTokError(400, "Invalid request body"))
	}

	// Return 5 dummy videos for immediate testing
	dummyVideos := []fiber.Map{
		{
			"Video": fiber.Map{
				"id":             1,
				"user_id":        1,
				"description":    "Welcome to Lomi Live! 🎉",
				"video":          "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/BigBuckBunny.mp4",
				"thumbnail":      "https://picsum.photos/400/600?random=1",
				"sound_id":       1,
				"view":           15000,
				"like":           1200,
				"comment_count":  45,
				"share":          30,
				"privacy_type":   "public",
				"allow_comments": 1,
				"allow_duet":     1,
				"created":        time.Now().Add(-24 * time.Hour).Format("2006-01-02 15:04:05"),
			},
			"User": fiber.Map{
				"id":          1,
				"username":    "lomi_official",
				"profile_pic": "https://picsum.photos/200?random=1",
				"verified":    1,
			},
			"Sound": fiber.Map{
				"id":    1,
				"title": "Original Sound",
				"sound": "",
			},
		},
		{
			"Video": fiber.Map{
				"id":             2,
				"user_id":        2,
				"description":    "Amazing dance moves! 💃",
				"video":          "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ElephantsDream.mp4",
				"thumbnail":      "https://picsum.photos/400/600?random=2",
				"sound_id":       2,
				"view":           25000,
				"like":           2100,
				"comment_count":  89,
				"share":          55,
				"privacy_type":   "public",
				"allow_comments": 1,
				"allow_duet":     1,
				"created":        time.Now().Add(-12 * time.Hour).Format("2006-01-02 15:04:05"),
			},
			"User": fiber.Map{
				"id":          2,
				"username":    "dancer_pro",
				"profile_pic": "https://picsum.photos/200?random=2",
				"verified":    0,
			},
			"Sound": fiber.Map{
				"id":    2,
				"title": "Trending Beat",
				"sound": "",
			},
		},
		{
			"Video": fiber.Map{
				"id":             3,
				"user_id":        3,
				"description":    "Cooking tutorial 🍳",
				"video":          "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerBlazes.mp4",
				"thumbnail":      "https://picsum.photos/400/600?random=3",
				"sound_id":       3,
				"view":           18000,
				"like":           1500,
				"comment_count":  67,
				"share":          42,
				"privacy_type":   "public",
				"allow_comments": 1,
				"allow_duet":     0,
				"created":        time.Now().Add(-6 * time.Hour).Format("2006-01-02 15:04:05"),
			},
			"User": fiber.Map{
				"id":          3,
				"username":    "chef_master",
				"profile_pic": "https://picsum.photos/200?random=3",
				"verified":    1,
			},
			"Sound": fiber.Map{
				"id":    3,
				"title": "Cooking Vibes",
				"sound": "",
			},
		},
		{
			"Video": fiber.Map{
				"id":             4,
				"user_id":        4,
				"description":    "Travel vlog ✈️",
				"video":          "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerEscapes.mp4",
				"thumbnail":      "https://picsum.photos/400/600?random=4",
				"sound_id":       4,
				"view":           32000,
				"like":           2800,
				"comment_count":  120,
				"share":          78,
				"privacy_type":   "public",
				"allow_comments": 1,
				"allow_duet":     1,
				"created":        time.Now().Add(-3 * time.Hour).Format("2006-01-02 15:04:05"),
			},
			"User": fiber.Map{
				"id":          4,
				"username":    "traveler_life",
				"profile_pic": "https://picsum.photos/200?random=4",
				"verified":    1,
			},
			"Sound": fiber.Map{
				"id":    4,
				"title": "Adventure Music",
				"sound": "",
			},
		},
		{
			"Video": fiber.Map{
				"id":             5,
				"user_id":        5,
				"description":    "Comedy sketch 😂",
				"video":          "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerFun.mp4",
				"thumbnail":      "https://picsum.photos/400/600?random=5",
				"sound_id":       5,
				"view":           45000,
				"like":           3900,
				"comment_count":  234,
				"share":          156,
				"privacy_type":   "public",
				"allow_comments": 1,
				"allow_duet":     1,
				"created":        time.Now().Add(-1 * time.Hour).Format("2006-01-02 15:04:05"),
			},
			"User": fiber.Map{
				"id":          5,
				"username":    "funny_guy",
				"profile_pic": "https://picsum.photos/200?random=5",
				"verified":    0,
			},
			"Sound": fiber.Map{
				"id":    5,
				"title": "Funny Sound",
				"sound": "",
			},
		},
	}

	// Paginate based on starting_point
	start := req.StartingPoint
	if start >= len(dummyVideos) {
		return c.JSON(tikTokSuccess([]fiber.Map{}))
	}

	end := start + 10
	if end > len(dummyVideos) {
		end = len(dummyVideos)
	}

	response := tikTokSuccess(dummyVideos[start:end])
	return c.JSON(response)
}

// ==================== 4. POST /api/liveStream ====================
// Start live streaming session
func (h *StreamingHandler) LiveStream(c *fiber.Ctx) error {
	var req struct {
		UserID    string `json:"user_id"`
		StartedAt string `json:"started_at"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.JSON(tikTokError(400, "Invalid request body"))
	}

	// Get user ID from context or request
	userID := c.Locals("user_id")
	if userID == nil && req.UserID != "" {
		parsedID, err := uuid.Parse(req.UserID)
		if err == nil {
			userID = parsedID
		}
	}

	if userID == nil {
		return c.JSON(tikTokError(401, "Unauthorized"))
	}

	// Generate unique streaming ID
	streamingID := generateStreamingID()

	// TODO: Create LiveStreaming record in database
	// For now, return the streaming ID which can be used with MediaMTX

	log.Printf("✅ Live stream started: user_id=%s, streaming_id=%s", userID, streamingID)

	response := tikTokSuccess(fiber.Map{
		"LiveStreaming": fiber.Map{
			"id":           streamingID,
			"user_id":      userID,
			"channel_name": streamingID,
			"started_at":   req.StartedAt,
			"status":       "live",
			// MediaMTX RTMP URL for streaming
			"rtmp_url": fmt.Sprintf("rtmp://localhost:1935/live/%s", streamingID),
			// HLS playback URL
			"playback_url": fmt.Sprintf("http://localhost:8888/live/%s/index.m3u8", streamingID),
		},
	})

	return c.JSON(response)
}

// ==================== 5. POST /api/sendGift ====================
// Send virtual gift during live stream or on video
func (h *StreamingHandler) SendGift(c *fiber.Ctx) error {
	var req struct {
		SenderID        string `json:"sender_id"`
		ReceiverID      string `json:"receiver_id"`
		VideoID         string `json:"video_id"`
		LiveStreamingID string `json:"live_streaming_id"`
		GiftID          string `json:"gift_id"`
		GiftCount       int    `json:"gift_count"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.JSON(tikTokError(400, "Invalid request body"))
	}

	// Get sender ID from context
	senderID := c.Locals("user_id")
	if senderID == nil {
		return c.JSON(tikTokError(401, "Unauthorized"))
	}

	// Parse receiver ID
	receiverID, err := uuid.Parse(req.ReceiverID)
	if err != nil {
		return c.JSON(tikTokError(400, "Invalid receiver_id"))
	}

	// Parse gift ID
	giftID, err := uuid.Parse(req.GiftID)
	if err != nil {
		return c.JSON(tikTokError(400, "Invalid gift_id"))
	}

	if req.GiftCount <= 0 {
		req.GiftCount = 1
	}

	var sender models.User
	var receiver models.User
	var gift models.Gift

	err = database.DB.Transaction(func(tx *gorm.DB) error {
		// Get sender
		if err := tx.Where("id = ?", senderID).First(&sender).Error; err != nil {
			return err
		}

		// Get receiver
		if err := tx.Where("id = ?", receiverID).First(&receiver).Error; err != nil {
			return err
		}

		// Get gift
		if err := tx.Where("id = ?", giftID).First(&gift).Error; err != nil {
			return err
		}

		// Calculate total cost
		totalCost := gift.CoinPrice * req.GiftCount

		// Check if sender has enough coins
		if sender.CoinBalance < totalCost {
			return fmt.Errorf("insufficient coins")
		}

		// Deduct coins from sender
		sender.CoinBalance -= totalCost
		sender.TotalSpent += totalCost
		if err := tx.Save(&sender).Error; err != nil {
			return err
		}

		// Add coins to receiver
		receiver.CoinBalance += totalCost
		receiver.TotalEarned += totalCost
		if err := tx.Save(&receiver).Error; err != nil {
			return err
		}

		// Create gift transaction
		giftTx := models.GiftTransaction{
			SenderID:   sender.ID,
			ReceiverID: receiver.ID,
			GiftID:     gift.ID,
			CoinAmount: totalCost,
			BirrValue:  gift.BirrValue * float64(req.GiftCount),
			GiftType:   gift.NameEn,
		}
		if err := tx.Create(&giftTx).Error; err != nil {
			return err
		}

		// Create coin transactions for both users
		senderCoinTx := models.CoinTransaction{
			UserID:            sender.ID,
			TransactionType:   models.TransactionTypeGiftSent,
			CoinAmount:        -totalCost,
			GiftTransactionID: &giftTx.ID,
			BalanceAfter:      sender.CoinBalance,
			PaymentStatus:     models.PaymentStatusCompleted,
		}
		if err := tx.Create(&senderCoinTx).Error; err != nil {
			return err
		}

		receiverCoinTx := models.CoinTransaction{
			UserID:            receiver.ID,
			TransactionType:   models.TransactionTypeGiftReceived,
			CoinAmount:        totalCost,
			GiftTransactionID: &giftTx.ID,
			BalanceAfter:      receiver.CoinBalance,
			PaymentStatus:     models.PaymentStatusCompleted,
		}
		if err := tx.Create(&receiverCoinTx).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Printf("❌ SendGift error: %v", err)
		if err.Error() == "insufficient coins" {
			return c.JSON(tikTokError(201, "You don't have sufficient coins"))
		}
		return c.JSON(tikTokError(500, "Could not send gift"))
	}

	log.Printf("✅ Gift sent: sender=%s, receiver=%s, gift=%s, count=%d", sender.ID, receiver.ID, gift.NameEn, req.GiftCount)

	response := tikTokSuccess(fiber.Map{
		"User": fiber.Map{
			"id":     sender.ID,
			"wallet": sender.CoinBalance,
		},
		"Gift": fiber.Map{
			"id":    gift.ID,
			"title": gift.NameEn,
			"coin":  gift.CoinPrice,
		},
	})

	return c.JSON(response)
}

// ==================== 6. POST /api/purchaseCoin ====================
// Purchase coins (in-app purchase)
func (h *StreamingHandler) PurchaseCoin(c *fiber.Ctx) error {
	var req struct {
		UserID        string `json:"user_id"`
		Coin          string `json:"coin"`
		Title         string `json:"title"`
		Price         string `json:"price"`
		TransactionID string `json:"transaction_id"`
		Device        string `json:"device"` // "ios" or "android"
	}

	if err := c.BodyParser(&req); err != nil {
		return c.JSON(tikTokError(400, "Invalid request body"))
	}

	// Get user ID from context
	userID := c.Locals("user_id")
	if userID == nil {
		return c.JSON(tikTokError(401, "Unauthorized"))
	}

	// Parse coin amount
	var coinAmount int
	fmt.Sscanf(req.Coin, "%d", &coinAmount)
	if coinAmount <= 0 {
		return c.JSON(tikTokError(400, "Invalid coin amount"))
	}

	// Parse price
	var price float64
	fmt.Sscanf(req.Price, "%f", &price)
	if price <= 0 {
		return c.JSON(tikTokError(400, "Invalid price"))
	}

	var user models.User

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Get user
		if err := tx.Where("id = ?", userID).First(&user).Error; err != nil {
			return err
		}

		// Add coins to user
		user.CoinBalance += coinAmount
		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		// Create coin transaction
		coinTx := models.CoinTransaction{
			UserID:           user.ID,
			TransactionType:  models.TransactionTypePurchase,
			CoinAmount:       coinAmount,
			BirrAmount:       price,
			PaymentMethod:    models.PaymentMethodTelebirr, // Default
			PaymentReference: req.TransactionID,
			PaymentStatus:    models.PaymentStatusCompleted,
			BalanceAfter:     user.CoinBalance,
			Metadata: models.JSONMap{
				"title":  req.Title,
				"device": req.Device,
			},
		}
		if err := tx.Create(&coinTx).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Printf("❌ PurchaseCoin error: %v", err)
		return c.JSON(tikTokError(500, "Could not process purchase"))
	}

	log.Printf("✅ Coins purchased: user=%s, coins=%d, price=%.2f", user.ID, coinAmount, price)

	response := tikTokSuccess(fiber.Map{
		"User": fiber.Map{
			"id":                   user.ID,
			"wallet":               user.CoinBalance,
			"total_all_time_coins": user.TotalEarned + coinAmount,
		},
		"Transaction": fiber.Map{
			"id":      999, // Dummy ID
			"coin":    coinAmount,
			"price":   price,
			"created": time.Now().Format("2006-01-02 15:04:05"),
		},
	})

	return c.JSON(response)
}

// ==================== Helper Functions ====================

func generateAuthToken(userID uuid.UUID) (string, error) {
	// For now, return a simple token
	// In production, use JWT with proper signing
	return fmt.Sprintf("auth_%s_%d", userID.String(), time.Now().Unix()), nil
}

func generateStreamingID() string {
	// Generate a random streaming ID
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
