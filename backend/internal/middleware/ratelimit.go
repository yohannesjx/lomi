package middleware

import (
	"lomi-backend/internal/database"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// RateLimitConfig configures rate limiting
type RateLimitConfig struct {
	MaxRequests int           // Maximum number of requests
	Window      time.Duration // Time window
	KeyPrefix   string        // Redis key prefix
}

// RateLimit creates a rate limiting middleware
func RateLimit(config RateLimitConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user ID from JWT token
		user := c.Locals("user")
		if user == nil {
			return c.Next()
		}

		token, ok := user.(*jwt.Token)
		if !ok {
			return c.Next()
		}

		claims := token.Claims.(jwt.MapClaims)
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			return c.Next()
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return c.Next()
		}

		// Create Redis key
		key := config.KeyPrefix + ":" + userID.String()
		windowSeconds := int(config.Window.Seconds())

		// Use Redis to track rate limits
		if database.RedisClient != nil {
			// Get current count
			countStr, err := database.RedisClient.Get(c.Context(), key).Result()
			if err != nil && err.Error() != "redis: nil" {
				// Redis error, allow request but log
				return c.Next()
			}

			count := 0
			if countStr != "" {
				count, _ = strconv.Atoi(countStr)
			}

			// Check if limit exceeded
			if count >= config.MaxRequests {
				return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
					"error": "Rate limit exceeded",
					"retry_after": windowSeconds,
				})
			}

			// Increment counter
			pipe := database.RedisClient.Pipeline()
			pipe.Incr(c.Context(), key)
			pipe.Expire(c.Context(), key, config.Window)
			_, err = pipe.Exec(c.Context())
			if err != nil {
				// Redis error, allow request
				return c.Next()
			}
		}

		return c.Next()
	}
}

// SwipeRateLimit limits swipes per hour
func SwipeRateLimit() fiber.Handler {
	return RateLimit(RateLimitConfig{
		MaxRequests: 100, // Max 100 swipes per hour
		Window:      time.Hour,
		KeyPrefix:   "ratelimit:swipe",
	})
}

// MessageRateLimit limits messages per minute
func MessageRateLimit() fiber.Handler {
	return RateLimit(RateLimitConfig{
		MaxRequests: 30, // Max 30 messages per minute
		Window:      time.Minute,
		KeyPrefix:   "ratelimit:message",
	})
}

// PurchaseRateLimit limits coin purchases per day
func PurchaseRateLimit() fiber.Handler {
	return RateLimit(RateLimitConfig{
		MaxRequests: 10, // Max 10 purchases per day
		Window:      24 * time.Hour,
		KeyPrefix:   "ratelimit:purchase",
	})
}

