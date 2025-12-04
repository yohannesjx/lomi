package handlers

import (
	"strconv"

	"lomi-backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type VideoHandler struct {
	videoService *services.VideoService
}

func NewVideoHandler(videoService *services.VideoService) *VideoHandler {
	return &VideoHandler{
		videoService: videoService,
	}
}

// ============================================
// VIDEO CONTENT ENDPOINTS
// ============================================

// ShowVideosAgainstUserID handles POST /api/v1/showVideosAgainstUserID
func (h *VideoHandler) ShowVideosAgainstUserID(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	viewerID := claims["user_id"].(string)

	// Get user_id from request body
	var reqBody struct {
		UserID string `json:"user_id"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "Invalid request body",
		})
	}

	userID := reqBody.UserID
	if userID == "" {
		userID = viewerID // If no user_id provided, show own videos
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	result, err := h.videoService.GetUserVideos(c.Context(), userID, viewerID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code": 500,
			"msg":  "Failed to get videos",
		})
	}

	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  "success",
		"data": result,
	})
}

// ShowUserLikedVideos handles POST /api/v1/showUserLikedVideos
func (h *VideoHandler) ShowUserLikedVideos(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	viewerID := claims["user_id"].(string)

	// Get user_id from request body
	var reqBody struct {
		UserID string `json:"user_id"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "Invalid request body",
		})
	}

	userID := reqBody.UserID
	if userID == "" {
		userID = viewerID // If no user_id provided, show own liked videos
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	result, err := h.videoService.GetUserLikedVideos(c.Context(), userID, viewerID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code": 500,
			"msg":  "Failed to get liked videos",
		})
	}

	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  "success",
		"data": result,
	})
}

// ShowUserRepostedVideos handles POST /api/v1/showUserRepostedVideos
func (h *VideoHandler) ShowUserRepostedVideos(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	viewerID := claims["user_id"].(string)

	// Get user_id from request body
	var reqBody struct {
		UserID string `json:"user_id"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "Invalid request body",
		})
	}

	userID := reqBody.UserID
	if userID == "" {
		userID = viewerID // If no user_id provided, show own reposts
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	result, err := h.videoService.GetUserRepostedVideos(c.Context(), userID, viewerID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code": 500,
			"msg":  "Failed to get reposted videos",
		})
	}

	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  "success",
		"data": result,
	})
}

// ShowFavouriteVideos handles POST /api/v1/showFavouriteVideos
func (h *VideoHandler) ShowFavouriteVideos(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	viewerID := claims["user_id"].(string)

	// Get user_id from request body
	var reqBody struct {
		UserID string `json:"user_id"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "Invalid request body",
		})
	}

	userID := reqBody.UserID
	if userID == "" {
		userID = viewerID // If no user_id provided, show own favorites
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	result, err := h.videoService.GetUserFavoriteVideos(c.Context(), userID, viewerID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code": 500,
			"msg":  "Failed to get favorite videos",
		})
	}

	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  "success",
		"data": result,
	})
}

// ============================================
// DRAFT VIDEOS
// ============================================

// ShowDraftVideos handles POST /api/v1/showDraftVideos
func (h *VideoHandler) ShowDraftVideos(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)
	
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	
	result, err := h.videoService.GetUserDraftVideos(c.Context(), userID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code": 500,
			"msg":  "Failed to get draft videos",
		})
	}
	
	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  "success",
		"data": result,
	})
}

// DeleteDraftVideo handles POST /api/v1/deleteDraftVideo
func (h *VideoHandler) DeleteDraftVideo(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)
	
	var reqBody struct {
		VideoID string `json:"video_id"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "Invalid request body",
		})
	}
	
	if reqBody.VideoID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  "video_id is required",
		})
	}
	
	err := h.videoService.DeleteDraftVideo(c.Context(), reqBody.VideoID, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code": 400,
			"msg":  err.Error(),
		})
	}
	
	return c.JSON(fiber.Map{
		"code": 200,
		"msg":  "Draft video deleted successfully",
	})
}
