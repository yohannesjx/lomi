package handlers

import (
	"fmt"
	"log"
	"lomi-backend/internal/database"
	"lomi-backend/internal/models"
	"lomi-backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Gift definitions matching the spec
var GiftCatalog = []struct {
	Type         string
	Name         string
	CoinPrice    int
	AnimationURL string
	SoundURL     string
}{
	{"rose", "Rose", 290, "/animations/rose.json", "/sounds/rose.mp3"},
	{"heart", "Heart", 499, "/animations/heart.json", "/sounds/heart.mp3"},
	{"diamond_ring", "Diamond Ring", 999, "/animations/diamond_ring.json", "/sounds/diamond_ring.mp3"},
	{"fireworks", "Fireworks", 1999, "/animations/fireworks.json", "/sounds/fireworks.mp3"},
	{"yacht", "Yacht", 4999, "/animations/yacht.json", "/sounds/yacht.mp3"},
	{"sports_car", "Sports Car", 9999, "/animations/sports_car.json", "/sounds/sports_car.mp3"},
	{"private_jet", "Private Jet", 29999, "/animations/private_jet.json", "/sounds/private_jet.mp3"},
	{"castle", "Castle", 79999, "/animations/castle.json", "/sounds/castle.mp3"},
	{"universe", "Universe", 149999, "/animations/universe.json", "/sounds/universe.mp3"},
	{"lomi_crown", "Lomi Crown", 299999, "/animations/lomi_crown.json", "/sounds/lomi_crown.mp3"},
}

// Coin purchase packs
var CoinPacks = []struct {
	ID       string
	Name     string
	ETBPrice float64
	Coins    int
}{
	{"spark", "Spark", 55, 600},
	{"flame", "Flame", 110, 1300},
	{"blaze", "Blaze", 275, 3500},
	{"inferno", "Inferno", 550, 8000},
	{"galaxy", "Galaxy", 1100, 18000},
	{"universe", "Universe", 5500, 100000},
}

// GetGiftShop returns all gifts with prices and animation URLs
func GetGiftShop(c *fiber.Ctx) error {
	gifts := make([]fiber.Map, 0, len(GiftCatalog))

	for _, gift := range GiftCatalog {
		// Calculate ETB value (1 LC = 0.1 ETB)
		etbValue := float64(gift.CoinPrice) * 0.1

		gifts = append(gifts, fiber.Map{
			"type":          gift.Type,
			"name":          gift.Name,
			"coin_price":    gift.CoinPrice,
			"etb_value":     etbValue,
			"animation_url": gift.AnimationURL,
			"sound_url":     gift.SoundURL,
		})
	}

	return c.JSON(fiber.Map{
		"gifts": gifts,
		"count": len(gifts),
	})
}

// GetWalletBalance returns user's current LC balance
func GetWalletBalance(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	userID, _ := uuid.Parse(userIDStr)

	var dbUser models.User
	if err := database.DB.First(&dbUser, "id = ?", userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(fiber.Map{
		"coin_balance": dbUser.CoinBalance,
		"total_spent":  dbUser.TotalSpent,
		"total_earned": dbUser.TotalEarned,
		"etb_value":    float64(dbUser.CoinBalance) * 0.1, // 1 LC = 0.1 ETB
	})
}

// BuyCoins initiates coin purchase (redirects to Telebirr/CBE Birr)
func BuyCoins(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	userID, _ := uuid.Parse(userIDStr)

	var req struct {
		PackID string `json:"pack_id" validate:"required"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Find pack
	var selectedPack *struct {
		ID       string
		Name     string
		ETBPrice float64
		Coins    int
	}
	for _, pack := range CoinPacks {
		if pack.ID == req.PackID {
			selectedPack = &pack
			break
		}
	}
	if selectedPack == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid pack ID"})
	}

	// Create pending transaction
	coinTx := models.CoinTransaction{
		UserID:          userID,
		TransactionType: models.TransactionTypePurchase,
		CoinAmount:      selectedPack.Coins,
		BirrAmount:      selectedPack.ETBPrice,
		PaymentMethod:   models.PaymentMethodTelebirr, // Default to Telebirr
		PaymentStatus:   models.PaymentStatusPending,
		BalanceAfter:    0, // Will be updated after payment
	}

	if err := database.DB.Create(&coinTx).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create transaction"})
	}

	// TODO: Generate Telebirr payment URL
	// For now, return payment URL structure
	paymentURL := "https://telebirr.com/pay?amount=" +
		fmt.Sprintf("%.2f", selectedPack.ETBPrice) +
		"&reference=" + coinTx.ID.String()

	return c.JSON(fiber.Map{
		"transaction_id": coinTx.ID,
		"pack_id":        selectedPack.ID,
		"pack_name":      selectedPack.Name,
		"etb_price":      selectedPack.ETBPrice,
		"coins":          selectedPack.Coins,
		"payment_url":    paymentURL,
		"webhook_url":    "/api/v1/wallet/buy/webhook", // Webhook endpoint
	})
}

// SendGiftLuxury sends a luxury gift (new implementation)
func SendGiftLuxury(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	senderID, _ := uuid.Parse(userIDStr)

	var req struct {
		ReceiverID string `json:"receiver_id" validate:"required"`
		GiftType   string `json:"gift_type" validate:"required"`
		MatchID    string `json:"match_id,omitempty"` // Optional: if sent in chat
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	receiverID, err := uuid.Parse(req.ReceiverID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid receiver ID"})
	}

	// Find gift in catalog
	var selectedGift *struct {
		Type         string
		Name         string
		CoinPrice    int
		AnimationURL string
		SoundURL     string
	}
	for _, gift := range GiftCatalog {
		if gift.Type == req.GiftType {
			selectedGift = &gift
			break
		}
	}
	if selectedGift == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid gift type"})
	}

	// Get sender
	var sender models.User
	if err := database.DB.First(&sender, "id = ?", senderID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Sender not found"})
	}

	// Check if sender has enough coins
	if sender.CoinBalance < selectedGift.CoinPrice {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":           "Insufficient coins",
			"required":        selectedGift.CoinPrice,
			"current_balance": sender.CoinBalance,
		})
	}

	// Get receiver
	var receiver models.User
	if err := database.DB.First(&receiver, "id = ?", receiverID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Receiver not found"})
	}

	// Calculate ETB value (1 LC = 0.1 ETB)
	etbValue := float64(selectedGift.CoinPrice) * 0.1

	// Start transaction
	tx := database.DB.Begin()

	// Deduct coins from sender
	sender.CoinBalance -= selectedGift.CoinPrice
	sender.TotalSpent += selectedGift.CoinPrice
	if err := tx.Save(&sender).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to deduct coins"})
	}

	// Add coins to receiver (they earn the full coin value)
	receiver.CoinBalance += selectedGift.CoinPrice
	receiver.TotalEarned += selectedGift.CoinPrice
	receiver.GiftBalance += etbValue
	if err := tx.Save(&receiver).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add coins to receiver"})
	}

	// Create gift transaction
	giftTransaction := models.GiftTransaction{
		SenderID:   senderID,
		ReceiverID: receiverID,
		GiftID:     uuid.Nil, // We don't have gift IDs in DB, using type instead
		CoinAmount: selectedGift.CoinPrice,
		BirrValue:  etbValue,
		GiftType:   selectedGift.Type,
	}

	if err := tx.Create(&giftTransaction).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create gift transaction"})
	}

	// Create coin transaction for sender
	coinTxSender := models.CoinTransaction{
		UserID:            senderID,
		TransactionType:   models.TransactionTypeGiftSent,
		CoinAmount:        -selectedGift.CoinPrice,
		BalanceAfter:      sender.CoinBalance,
		GiftTransactionID: &giftTransaction.ID,
	}
	if err := tx.Create(&coinTxSender).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create sender coin transaction"})
	}

	// Create coin transaction for receiver
	coinTxReceiver := models.CoinTransaction{
		UserID:            receiverID,
		TransactionType:   models.TransactionTypeGiftReceived,
		CoinAmount:        selectedGift.CoinPrice,
		BalanceAfter:      receiver.CoinBalance,
		GiftTransactionID: &giftTransaction.ID,
	}
	if err := tx.Create(&coinTxReceiver).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create receiver coin transaction"})
	}

	// If sent in chat, create a message
	if req.MatchID != "" {
		matchID, _ := uuid.Parse(req.MatchID)
		message := models.Message{
			MatchID:     &matchID,
			SenderID:    senderID,
			ReceiverID:  &receiverID,
			MessageType: models.MessageTypeGift,
			IsLive:      false,
		}
		if err := tx.Create(&message).Error; err == nil {
			giftTransaction.MessageID = &message.ID
			tx.Save(&giftTransaction)
		}
	}

	tx.Commit()

	// Send push notification (async)
	go func() {
		if services.NotificationSvc != nil {
			services.NotificationSvc.NotifyGiftReceived(giftTransaction, sender)
		}
	}()

	// If gift is 29,999+ coins, send push to all users
	if selectedGift.CoinPrice >= 29999 {
		go func() {
			// TODO: Implement broadcast notification
			log.Printf("🎉 BIG GIFT ALERT: %s sent a %s (%d LC) to %s in %s!",
				sender.Name, selectedGift.Name, selectedGift.CoinPrice, receiver.Name, receiver.City)
		}()
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Gift sent successfully",
		"gift": fiber.Map{
			"type":          selectedGift.Type,
			"name":          selectedGift.Name,
			"coin_price":    selectedGift.CoinPrice,
			"animation_url": selectedGift.AnimationURL,
			"sound_url":     selectedGift.SoundURL,
		},
		"sender_balance":   sender.CoinBalance,
		"receiver_balance": receiver.CoinBalance,
	})
}

// GetGiftsReceived returns list of gifts user received (for cashout page)
func GetGiftsReceived(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	userID, _ := uuid.Parse(userIDStr)

	var gifts []models.GiftTransaction
	if err := database.DB.Where("receiver_id = ?", userID).
		Order("created_at DESC").
		Limit(100).
		Find(&gifts).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch gifts"})
	}

	// Format response
	giftsList := make([]fiber.Map, 0, len(gifts))
	totalCoins := 0
	totalETB := 0.0

	for _, gift := range gifts {
		var sender models.User
		database.DB.First(&sender, "id = ?", gift.SenderID)

		giftsList = append(giftsList, fiber.Map{
			"id":          gift.ID,
			"sender_name": sender.Name,
			"gift_type":   gift.GiftType,
			"coins":       gift.CoinAmount,
			"etb_value":   gift.BirrValue,
			"sent_at":     gift.CreatedAt,
		})

		totalCoins += gift.CoinAmount
		totalETB += gift.BirrValue
	}

	return c.JSON(fiber.Map{
		"gifts":       giftsList,
		"total_coins": totalCoins,
		"total_etb":   totalETB,
		"count":       len(gifts),
	})
}

// RequestCashout creates a cashout request (min 50,000 LC)
func RequestCashout(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	userID, _ := uuid.Parse(userIDStr)

	var req struct {
		Coins          int    `json:"coins" validate:"required,min=50000"`
		PaymentMethod  string `json:"payment_method" validate:"required"`
		PaymentAccount string `json:"payment_account" validate:"required"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if req.Coins < 50000 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Minimum cashout is 50,000 LC",
		})
	}

	// Get user
	var dbUser models.User
	if err := database.DB.First(&dbUser, "id = ?", userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Check if user has enough coins
	if dbUser.CoinBalance < req.Coins {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":           "Insufficient coins",
			"required":        req.Coins,
			"current_balance": dbUser.CoinBalance,
		})
	}

	// Calculate ETB amount (1 LC = 0.1 ETB)
	etbAmount := float64(req.Coins) * 0.1
	platformFeePercentage := 25 // 25% platform fee
	platformFeeAmount := etbAmount * float64(platformFeePercentage) / 100.0
	netAmount := etbAmount - platformFeeAmount

	// Parse payment method
	var paymentMethod models.PaymentMethod
	switch req.PaymentMethod {
	case "telebirr":
		paymentMethod = models.PaymentMethodTelebirr
	case "cbe_birr":
		paymentMethod = models.PaymentMethodCbeBirr
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payment method"})
	}

	// Create payout request
	payout := models.Payout{
		UserID:                userID,
		Coins:                 req.Coins,
		GiftBalanceAmount:     etbAmount,
		PlatformFeePercentage: platformFeePercentage,
		PlatformFeeAmount:     platformFeeAmount,
		NetAmount:             netAmount,
		PaymentMethod:         paymentMethod,
		PaymentAccount:        req.PaymentAccount,
		Status:                models.PayoutStatusPending,
	}

	if err := database.DB.Create(&payout).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create cashout request"})
	}

	// Deduct coins from user (they'll get it back if rejected)
	dbUser.CoinBalance -= req.Coins
	if err := database.DB.Save(&dbUser).Error; err != nil {
		log.Printf("⚠️ Failed to deduct coins for cashout: %v", err)
		// Don't fail the request, admin can handle it
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Cashout request created",
		"payout": fiber.Map{
			"id":              payout.ID,
			"coins":           payout.Coins,
			"etb_amount":      etbAmount,
			"platform_fee":    platformFeeAmount,
			"net_amount":      netAmount,
			"payment_method":  payout.PaymentMethod,
			"payment_account": payout.PaymentAccount,
			"status":          payout.Status,
			"created_at":      payout.CreatedAt,
		},
	})
}
