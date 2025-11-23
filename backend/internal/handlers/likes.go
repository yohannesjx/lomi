package handlers

import (
	"lomi-backend/internal/database"
	"lomi-backend/internal/models"
	"lomi-backend/internal/services"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// GetPendingLikes returns users who liked the current user but haven't been liked back
func GetPendingLikes(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	userID, _ := uuid.Parse(userIDStr)

	// Get users who liked me (swiped on me with like/super_like)
	var swipes []models.Swipe
	database.DB.Where("swiped_id = ? AND action IN ?", userID, []models.SwipeAction{models.SwipeActionLike, models.SwipeActionSuperLike}).
		Order("created_at DESC").
		Find(&swipes)

	// Check if I've liked them back (if so, they're a match, not a pending like)
	var myLikes []uuid.UUID
	database.DB.Model(&models.Swipe{}).
		Where("swiper_id = ? AND swiped_id IN (SELECT swiper_id FROM swipes WHERE swiped_id = ? AND action IN ?)", userID, userID, []models.SwipeAction{models.SwipeActionLike, models.SwipeActionSuperLike}).
		Pluck("swiped_id", &myLikes)

	// Filter out users I've already liked back
	pendingLikerIDs := make([]uuid.UUID, 0)
	for _, swipe := range swipes {
		alreadyLiked := false
		for _, likedID := range myLikes {
			if swipe.SwiperID == likedID {
				alreadyLiked = true
				break
			}
		}
		if !alreadyLiked {
			pendingLikerIDs = append(pendingLikerIDs, swipe.SwiperID)
		}
	}

	// Get user details for pending likers
	var pendingUsers []models.User
	if len(pendingLikerIDs) > 0 {
		database.DB.Where("id IN ? AND is_active = ?", pendingLikerIDs, true).Find(&pendingUsers)
	}

	// Get swipe timestamps for each user
	type PendingLike struct {
		User       models.User `json:"user"`
		LikedAt    time.Time   `json:"liked_at"`
		IsRevealed bool        `json:"is_revealed"` // For future: track if user has revealed this person
	}

	pendingLikes := make([]PendingLike, 0)
	for _, u := range pendingUsers {
		// Find the swipe timestamp
		var swipe models.Swipe
		database.DB.Where("swiper_id = ? AND swiped_id = ?", u.ID, userID).
			Order("created_at DESC").
			First(&swipe)

		pendingLikes = append(pendingLikes, PendingLike{
			User:       u,
			LikedAt:    swipe.CreatedAt,
			IsRevealed: false, // TODO: Track reveals in a separate table if needed
		})
	}

	// Get current user's daily free reveal status
	var currentUser models.User
	database.DB.First(&currentUser, "id = ?", userID)

	// Check if free reveal resets (Addis time is UTC+3)
	now := time.Now()
	addisTime := now.UTC().Add(3 * time.Hour)
	today := time.Date(addisTime.Year(), addisTime.Month(), addisTime.Day(), 0, 0, 0, 0, time.UTC)

	hasFreeReveal := false
	if currentUser.LastRevealDate.IsZero() || currentUser.LastRevealDate.Before(today) {
		hasFreeReveal = true
	} else {
		hasFreeReveal = !currentUser.DailyFreeRevealUsed
	}

	// Calculate reset time (midnight Addis time)
	resetTime := today.Add(24 * time.Hour)
	if addisTime.Hour() >= 0 && addisTime.Hour() < 3 {
		resetTime = today // If before 3 AM UTC, reset is today
	}

	return c.JSON(fiber.Map{
		"pending_likes":   pendingLikes,
		"count":           len(pendingLikes),
		"has_free_reveal": hasFreeReveal,
		"reset_at":        resetTime,
	})
}

// RevealLike handles revealing a pending like (free or paid)
func RevealLike(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	userID, _ := uuid.Parse(userIDStr)

	var req struct {
		RevealAll bool   `json:"reveal_all"` // If true, reveal all for 299 coins
		TargetID  string `json:"target_id"`  // If reveal_all is false, reveal this specific user for 99 coins
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	var currentUser models.User
	if err := database.DB.First(&currentUser, "id = ?", userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Get pending likes
	var swipes []models.Swipe
	database.DB.Where("swiped_id = ? AND action IN ?", userID, []models.SwipeAction{models.SwipeActionLike, models.SwipeActionSuperLike}).
		Order("created_at DESC").
		Find(&swipes)

	// Filter out users already liked back
	var myLikes []uuid.UUID
	database.DB.Model(&models.Swipe{}).
		Where("swiper_id = ?", userID).
		Pluck("swiped_id", &myLikes)

	pendingLikerIDs := make([]uuid.UUID, 0)
	for _, swipe := range swipes {
		alreadyLiked := false
		for _, likedID := range myLikes {
			if swipe.SwiperID == likedID {
				alreadyLiked = true
				break
			}
		}
		if !alreadyLiked {
			pendingLikerIDs = append(pendingLikerIDs, swipe.SwiperID)
		}
	}

	if len(pendingLikerIDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No pending likes to reveal"})
	}

	// Check daily free reveal status
	now := time.Now()
	addisTime := now.UTC().Add(3 * time.Hour)
	today := time.Date(addisTime.Year(), addisTime.Month(), addisTime.Day(), 0, 0, 0, 0, time.UTC)

	hasFreeReveal := false
	if currentUser.LastRevealDate.IsZero() || currentUser.LastRevealDate.Before(today) {
		hasFreeReveal = true
	} else {
		hasFreeReveal = !currentUser.DailyFreeRevealUsed
	}

	// Reset daily free reveal if it's a new day
	if currentUser.LastRevealDate.IsZero() || currentUser.LastRevealDate.Before(today) {
		currentUser.DailyFreeRevealUsed = false
		currentUser.LastRevealDate = today
	}

	var cost int
	var revealedIDs []uuid.UUID

	if req.RevealAll {
		// Reveal all for 299 coins (or free if first reveal of day)
		if hasFreeReveal {
			// First reveal is free, rest cost 299
			if len(pendingLikerIDs) > 1 {
				cost = 299
			} else {
				cost = 0 // Only one like, use free reveal
			}
			revealedIDs = pendingLikerIDs
		} else {
			cost = 299
			revealedIDs = pendingLikerIDs
		}
	} else {
		// Reveal one for 99 coins (or free if first reveal of day)
		if hasFreeReveal {
			cost = 0
		} else {
			cost = 99
		}

		if req.TargetID != "" {
			targetID, err := uuid.Parse(req.TargetID)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid target ID"})
			}
			// Verify target is in pending likes
			found := false
			for _, id := range pendingLikerIDs {
				if id == targetID {
					found = true
					revealedIDs = []uuid.UUID{targetID}
					break
				}
			}
			if !found {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Target user not in pending likes"})
			}
		} else {
			// Reveal random one
			if len(pendingLikerIDs) > 0 {
				revealedIDs = []uuid.UUID{pendingLikerIDs[0]} // Take first (most recent)
			}
		}
	}

	// Check coin balance
	if cost > 0 && currentUser.CoinBalance < cost {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":    "Insufficient coins",
			"required": cost,
			"balance":  currentUser.CoinBalance,
		})
	}

	// Deduct coins
	if cost > 0 {
		currentUser.CoinBalance -= cost

		// Create coin transaction record
		transaction := models.CoinTransaction{
			UserID:          userID,
			TransactionType: models.TransactionTypeReveal,
			CoinAmount:      -cost,
			BalanceAfter:    currentUser.CoinBalance,
			Metadata: models.JSONMap{
				"reveal_type": map[string]interface{}{
					"reveal_all":     req.RevealAll,
					"revealed_count": len(revealedIDs),
				},
			},
		}
		database.DB.Create(&transaction)
	}

	// Update daily free reveal status
	if hasFreeReveal && cost == 0 {
		currentUser.DailyFreeRevealUsed = true
		currentUser.LastRevealDate = today
	}

	database.DB.Save(&currentUser)

	// Get revealed users
	var revealedUsers []models.User
	if len(revealedIDs) > 0 {
		database.DB.Where("id IN ?", revealedIDs).Find(&revealedUsers)
	}

	// Send push notifications to revealed users
	go func() {
		if services.NotificationSvc != nil {
			for _, revealedUser := range revealedUsers {
				services.NotificationSvc.NotifySomeoneViewedProfile(revealedUser.ID, userID)
			}
		}
	}()

	return c.JSON(fiber.Map{
		"revealed_users":  revealedUsers,
		"coins_deducted":  cost,
		"new_balance":     currentUser.CoinBalance,
		"has_free_reveal": false, // After reveal, no more free reveal today
	})
}
