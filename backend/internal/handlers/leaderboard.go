package handlers

import (
	"lomi-backend/internal/database"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GetTopGiftedUsers returns leaderboard of users who received the most gifts
func GetTopGiftedUsers(c *fiber.Ctx) error {
	timeframe := c.Query("timeframe", "week") // week, month, all
	limit := c.QueryInt("limit", 20)

	var startDate time.Time
	now := time.Now()

	switch timeframe {
	case "week":
		startDate = now.AddDate(0, 0, -7)
	case "month":
		startDate = now.AddDate(0, -1, 0)
	case "all":
		startDate = time.Time{} // No filter
	default:
		startDate = now.AddDate(0, 0, -7)
	}

	// Query to get top users by gift balance received
	type LeaderboardEntry struct {
		UserID   string  `json:"user_id"`
		Name     string  `json:"name"`
		Avatar   string  `json:"avatar"`
		Amount   float64 `json:"amount"`
		Rank     int     `json:"rank"`
	}

	var results []LeaderboardEntry

	query := database.DB.Table("gift_transactions").
		Select(`
			users.id as user_id,
			users.name,
			COALESCE(media.url, '') as avatar,
			SUM(gift_transactions.birr_value) as amount
		`).
		Joins("JOIN users ON gift_transactions.receiver_id = users.id").
		Joins("LEFT JOIN media ON users.id = media.user_id AND media.media_type = 'photo' AND media.display_order = 1 AND media.is_approved = true").
		Group("users.id, users.name, media.url")

	if !startDate.IsZero() {
		query = query.Where("gift_transactions.created_at >= ?", startDate)
	}

	query = query.Order("amount DESC").Limit(limit)

	rows, err := query.Rows()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch leaderboard"})
	}
	defer rows.Close()

	rank := 1
	for rows.Next() {
		var entry LeaderboardEntry
		var avatar *string
		if err := rows.Scan(&entry.UserID, &entry.Name, &avatar, &entry.Amount); err != nil {
			continue
		}
		if avatar != nil {
			entry.Avatar = *avatar
		}
		entry.Rank = rank
		results = append(results, entry)
		rank++
	}

	return c.JSON(fiber.Map{
		"leaderboard": results,
		"timeframe":   timeframe,
		"count":        len(results),
	})
}

