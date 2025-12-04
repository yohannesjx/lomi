package services

import (
	"context"

	"lomi-backend/internal/models"
	"lomi-backend/internal/repositories"
)

type VideoService struct {
	videoRepo *repositories.VideoRepository
}

func NewVideoService(videoRepo *repositories.VideoRepository) *VideoService {
	return &VideoService{
		videoRepo: videoRepo,
	}
}

// ============================================
// VIDEO CONTENT QUERIES
// ============================================

// GetUserVideos gets videos posted by a user
func (s *VideoService) GetUserVideos(ctx context.Context, userID, viewerID string, page, limit int) (*models.VideoListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 20
	}

	videos, totalCount, err := s.videoRepo.GetUserVideos(ctx, userID, viewerID, page, limit)
	if err != nil {
		return nil, err
	}

	hasMore := int64((page * limit)) < totalCount

	return &models.VideoListResponse{
		Videos:     videos,
		TotalCount: totalCount,
		Page:       page,
		Limit:      limit,
		HasMore:    hasMore,
	}, nil
}

// GetUserLikedVideos gets videos liked by a user
func (s *VideoService) GetUserLikedVideos(ctx context.Context, userID, viewerID string, page, limit int) (*models.VideoListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 20
	}

	videos, totalCount, err := s.videoRepo.GetUserLikedVideos(ctx, userID, viewerID, page, limit)
	if err != nil {
		return nil, err
	}

	hasMore := int64((page * limit)) < totalCount

	return &models.VideoListResponse{
		Videos:     videos,
		TotalCount: totalCount,
		Page:       page,
		Limit:      limit,
		HasMore:    hasMore,
	}, nil
}

// GetUserRepostedVideos gets videos reposted by a user
func (s *VideoService) GetUserRepostedVideos(ctx context.Context, userID, viewerID string, page, limit int) (*models.VideoListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 20
	}

	videos, totalCount, err := s.videoRepo.GetUserRepostedVideos(ctx, userID, viewerID, page, limit)
	if err != nil {
		return nil, err
	}

	hasMore := int64((page * limit)) < totalCount

	return &models.VideoListResponse{
		Videos:     videos,
		TotalCount: totalCount,
		Page:       page,
		Limit:      limit,
		HasMore:    hasMore,
	}, nil
}

// GetUserFavoriteVideos gets videos favorited by a user
func (s *VideoService) GetUserFavoriteVideos(ctx context.Context, userID, viewerID string, page, limit int) (*models.VideoListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 20
	}

	videos, totalCount, err := s.videoRepo.GetUserFavoriteVideos(ctx, userID, viewerID, page, limit)
	if err != nil {
		return nil, err
	}

	hasMore := int64((page * limit)) < totalCount

	return &models.VideoListResponse{
		Videos:     videos,
		TotalCount: totalCount,
		Page:       page,
		Limit:      limit,
		HasMore:    hasMore,
	}, nil
}
