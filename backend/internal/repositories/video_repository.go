package repositories

import (
	"context"
	"fmt"

	"lomi-backend/internal/models"

	"github.com/jmoiron/sqlx"
)

type VideoRepository struct {
	db *sqlx.DB
}

func NewVideoRepository(db *sqlx.DB) *VideoRepository {
	return &VideoRepository{db: db}
}

// ============================================
// VIDEO QUERIES
// ============================================

// GetUserVideos gets all videos posted by a user
func (r *VideoRepository) GetUserVideos(ctx context.Context, userID, viewerID string, page, limit int) ([]models.VideoResponse, int64, error) {
	offset := (page - 1) * limit

	// Get total count
	var totalCount int64
	countQuery := `
		SELECT COUNT(*) 
		FROM videos 
		WHERE user_id = $1::uuid AND deleted_at IS NULL
	`
	err := r.db.GetContext(ctx, &totalCount, countQuery, userID)
	if err != nil {
		return nil, 0, err
	}

	// Get videos with user info and engagement status
	query := `
		SELECT 
			v.*,
			u.name as user_name,
			u.is_verified as user_is_verified,
			EXISTS(SELECT 1 FROM video_likes WHERE video_id = v.id AND user_id = $2::uuid) as is_liked,
			EXISTS(SELECT 1 FROM video_favorites WHERE video_id = v.id AND user_id = $2::uuid) as is_favorited,
			EXISTS(SELECT 1 FROM video_reposts WHERE video_id = v.id AND user_id = $2::uuid) as is_reposted,
			EXISTS(SELECT 1 FROM follows WHERE follower_id = $2::uuid AND following_id = v.user_id) as is_following
		FROM videos v
		JOIN users u ON v.user_id = u.id
		WHERE v.user_id = $1::uuid AND v.deleted_at IS NULL
		ORDER BY v.created_at DESC
		LIMIT $3 OFFSET $4
	`

	videos := make([]models.VideoResponse, 0)
	err = r.db.SelectContext(ctx, &videos, query, userID, viewerID, limit, offset)
	if err != nil {
		return videos, totalCount, err
	}

	return videos, totalCount, nil
}

// GetUserLikedVideos gets all videos liked by a user
func (r *VideoRepository) GetUserLikedVideos(ctx context.Context, userID, viewerID string, page, limit int) ([]models.VideoResponse, int64, error) {
	offset := (page - 1) * limit

	// Get total count
	var totalCount int64
	countQuery := `
		SELECT COUNT(*) 
		FROM video_likes vl
		JOIN videos v ON vl.video_id = v.id
		WHERE vl.user_id = $1::uuid AND v.deleted_at IS NULL
	`
	err := r.db.GetContext(ctx, &totalCount, countQuery, userID)
	if err != nil {
		return nil, 0, err
	}

	// Get liked videos
	query := `
		SELECT 
			v.*,
			u.name as user_name,
			u.is_verified as user_is_verified,
			true as is_liked,
			EXISTS(SELECT 1 FROM video_favorites WHERE video_id = v.id AND user_id = $2::uuid) as is_favorited,
			EXISTS(SELECT 1 FROM video_reposts WHERE video_id = v.id AND user_id = $2::uuid) as is_reposted,
			EXISTS(SELECT 1 FROM follows WHERE follower_id = $2::uuid AND following_id = v.user_id) as is_following
		FROM video_likes vl
		JOIN videos v ON vl.video_id = v.id
		JOIN users u ON v.user_id = u.id
		WHERE vl.user_id = $1::uuid AND v.deleted_at IS NULL
		ORDER BY vl.created_at DESC
		LIMIT $3 OFFSET $4
	`

	videos := make([]models.VideoResponse, 0)
	err = r.db.SelectContext(ctx, &videos, query, userID, viewerID, limit, offset)
	if err != nil {
		return videos, totalCount, err
	}

	return videos, totalCount, nil
}

// GetUserRepostedVideos gets all videos reposted by a user
func (r *VideoRepository) GetUserRepostedVideos(ctx context.Context, userID, viewerID string, page, limit int) ([]models.VideoResponse, int64, error) {
	offset := (page - 1) * limit

	// Get total count
	var totalCount int64
	countQuery := `
		SELECT COUNT(*) 
		FROM video_reposts vr
		JOIN videos v ON vr.video_id = v.id
		WHERE vr.user_id = $1::uuid AND v.deleted_at IS NULL
	`
	err := r.db.GetContext(ctx, &totalCount, countQuery, userID)
	if err != nil {
		return nil, 0, err
	}

	// Get reposted videos
	query := `
		SELECT 
			v.*,
			u.name as user_name,
			u.is_verified as user_is_verified,
			EXISTS(SELECT 1 FROM video_likes WHERE video_id = v.id AND user_id = $2::uuid) as is_liked,
			EXISTS(SELECT 1 FROM video_favorites WHERE video_id = v.id AND user_id = $2::uuid) as is_favorited,
			true as is_reposted,
			EXISTS(SELECT 1 FROM follows WHERE follower_id = $2::uuid AND following_id = v.user_id) as is_following
		FROM video_reposts vr
		JOIN videos v ON vr.video_id = v.id
		JOIN users u ON v.user_id = u.id
		WHERE vr.user_id = $1::uuid AND v.deleted_at IS NULL
		ORDER BY vr.created_at DESC
		LIMIT $3 OFFSET $4
	`

	videos := make([]models.VideoResponse, 0)
	err = r.db.SelectContext(ctx, &videos, query, userID, viewerID, limit, offset)
	if err != nil {
		return videos, totalCount, err
	}

	return videos, totalCount, nil
}

// GetUserFavoriteVideos gets all videos favorited by a user
func (r *VideoRepository) GetUserFavoriteVideos(ctx context.Context, userID, viewerID string, page, limit int) ([]models.VideoResponse, int64, error) {
	offset := (page - 1) * limit

	// Get total count
	var totalCount int64
	countQuery := `
		SELECT COUNT(*) 
		FROM video_favorites vf
		JOIN videos v ON vf.video_id = v.id
		WHERE vf.user_id = $1::uuid AND v.deleted_at IS NULL
	`
	err := r.db.GetContext(ctx, &totalCount, countQuery, userID)
	if err != nil {
		return nil, 0, err
	}

	// Get favorite videos
	query := `
		SELECT 
			v.*,
			u.name as user_name,
			u.is_verified as user_is_verified,
			EXISTS(SELECT 1 FROM video_likes WHERE video_id = v.id AND user_id = $2::uuid) as is_liked,
			true as is_favorited,
			EXISTS(SELECT 1 FROM video_reposts WHERE video_id = v.id AND user_id = $2::uuid) as is_reposted,
			EXISTS(SELECT 1 FROM follows WHERE follower_id = $2::uuid AND following_id = v.user_id) as is_following
		FROM video_favorites vf
		JOIN videos v ON vf.video_id = v.id
		JOIN users u ON v.user_id = u.id
		WHERE vf.user_id = $1::uuid AND v.deleted_at IS NULL
		ORDER BY vf.created_at DESC
		LIMIT $3 OFFSET $4
	`

	videos := make([]models.VideoResponse, 0)
	err = r.db.SelectContext(ctx, &videos, query, userID, viewerID, limit, offset)
	if err != nil {
		return videos, totalCount, err
	}

	return videos, totalCount, nil
}

// ============================================
// DRAFT VIDEOS
// ============================================

// GetUserDraftVideos gets all draft videos for a user
func (r *VideoRepository) GetUserDraftVideos(ctx context.Context, userID string, page, limit int) ([]models.Video, int64, error) {
	offset := (page - 1) * limit

	// Get total count
	var totalCount int64
	countQuery := `
		SELECT COUNT(*) 
		FROM videos 
		WHERE user_id = $1::uuid AND is_draft = TRUE AND deleted_at IS NULL
	`
	err := r.db.GetContext(ctx, &totalCount, countQuery, userID)
	if err != nil {
		return nil, 0, err
	}

	// Get draft videos
	query := `
		SELECT *
		FROM videos
		WHERE user_id = $1::uuid AND is_draft = TRUE AND deleted_at IS NULL
		ORDER BY draft_saved_at DESC, created_at DESC
		LIMIT $2 OFFSET $3
	`

	videos := make([]models.Video, 0)
	err = r.db.SelectContext(ctx, &videos, query, userID, limit, offset)
	if err != nil {
		return videos, totalCount, err
	}

	return videos, totalCount, nil
}

// DeleteDraftVideo deletes a draft video
func (r *VideoRepository) DeleteDraftVideo(ctx context.Context, videoID, userID string) error {
	query := `
		UPDATE videos 
		SET deleted_at = NOW() 
		WHERE id = $1::uuid AND user_id = $2::uuid AND is_draft = TRUE AND deleted_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, query, videoID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("draft video not found or already deleted")
	}

	return nil
}
