package handlers

import (
	"lomi-backend/internal/database"
	"lomi-backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type UpdateProfileRequest struct {
	Name             string                 `json:"name"`
	Age              int                    `json:"age"`
	Gender           string                 `json:"gender"`
	Bio              string                 `json:"bio"`
	City             string                 `json:"city"`
	Interests        []string               `json:"interests"`
	RelationshipGoal string                 `json:"relationship_goal"`
	Preferences      map[string]interface{} `json:"preferences"`
}

func GetMe(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	var dbUser models.User
	if err := database.DB.First(&dbUser, "id = ?", userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(dbUser)
}

func UpdateProfile(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	var req UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var dbUser models.User
	if err := database.DB.First(&dbUser, "id = ?", userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Update fields
	if req.Name != "" {
		dbUser.Name = req.Name
	}
	if req.Age > 0 {
		dbUser.Age = req.Age
	}
	if req.Gender != "" {
		dbUser.Gender = models.Gender(req.Gender)
	}
	if req.Bio != "" {
		dbUser.Bio = req.Bio
	}
	if req.City != "" {
		dbUser.City = req.City
	}
	if len(req.Interests) > 0 {
		dbUser.Interests = req.Interests
	}
	if req.RelationshipGoal != "" {
		dbUser.RelationshipGoal = models.RelationshipGoal(req.RelationshipGoal)
	}
	if req.Preferences != nil {
		// Merge with existing preferences
		if dbUser.Preferences == nil {
			dbUser.Preferences = make(models.JSONMap)
		}
		for key, value := range req.Preferences {
			dbUser.Preferences[key] = value
		}
	}

	// Profile is considered complete if basic info is present (City check is done elsewhere)

	if err := database.DB.Save(&dbUser).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update profile"})
	}

	return c.JSON(dbUser)
}
