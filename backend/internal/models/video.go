package models

import "time"

// ============================================
// SOCIAL CONTENT MODELS
// ============================================

// Video represents a user-generated video
type Video struct {
	ID               string     `json:"id" db:"id"`
	UserID           string     `json:"user_id" db:"user_id"`
	VideoURL         string     `json:"video_url" db:"video_url"`
	ThumbnailURL     *string    `json:"thumbnail_url,omitempty" db:"thumbnail_url"`
	DurationSeconds  int        `json:"duration_seconds" db:"duration_seconds"`
	Title            *string    `json:"title,omitempty" db:"title"`
	Description      *string    `json:"description,omitempty" db:"description"`
	Hashtags         []string   `json:"hashtags" db:"hashtags"`
	IsPrivate        bool       `json:"is_private" db:"is_private"`
	AllowComments    bool       `json:"allow_comments" db:"allow_comments"`
	AllowDuet        bool       `json:"allow_duet" db:"allow_duet"`
	AllowStitch      bool       `json:"allow_stitch" db:"allow_stitch"`
	ViewsCount       int        `json:"views_count" db:"views_count"`
	LikesCount       int        `json:"likes_count" db:"likes_count"`
	CommentsCount    int        `json:"comments_count" db:"comments_count"`
	SharesCount      int        `json:"shares_count" db:"shares_count"`
	IsApproved       bool       `json:"is_approved" db:"is_approved"`
	ModerationStatus string     `json:"moderation_status" db:"moderation_status"`
	ModerationNotes  *string    `json:"moderation_notes,omitempty" db:"moderation_notes"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// VideoLike represents a like on a video
type VideoLike struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	VideoID   string    `json:"video_id" db:"video_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// VideoRepost represents a repost/share of a video
type VideoRepost struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	VideoID   string    `json:"video_id" db:"video_id"`
	Caption   *string   `json:"caption,omitempty" db:"caption"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// VideoFavorite represents a favorited video
type VideoFavorite struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	VideoID   string    `json:"video_id" db:"video_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// VideoComment represents a comment on a video
type VideoComment struct {
	ID              string     `json:"id" db:"id"`
	UserID          string     `json:"user_id" db:"user_id"`
	VideoID         string     `json:"video_id" db:"video_id"`
	ParentCommentID *string    `json:"parent_comment_id,omitempty" db:"parent_comment_id"`
	Content         string     `json:"content" db:"content"`
	LikesCount      int        `json:"likes_count" db:"likes_count"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// ============================================
// REQUEST/RESPONSE DTOs
// ============================================

// VideoResponse represents a video with user info and engagement status
type VideoResponse struct {
	Video
	// User info
	UserName       string  `json:"user_name" db:"user_name"`
	UserProfilePic *string `json:"user_profile_pic,omitempty" db:"user_profile_pic"`
	UserIsVerified bool    `json:"user_is_verified" db:"user_is_verified"`

	// Engagement status for current user
	IsLiked     bool `json:"is_liked" db:"is_liked"`
	IsFavorited bool `json:"is_favorited" db:"is_favorited"`
	IsReposted  bool `json:"is_reposted" db:"is_reposted"`
	IsFollowing bool `json:"is_following" db:"is_following"`
}

// GetUserVideosRequest represents request to get user's videos
type GetUserVideosRequest struct {
	UserID string `json:"user_id" validate:"required"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
}

// VideoListResponse represents a paginated list of videos
type VideoListResponse struct {
	Videos     []VideoResponse `json:"videos"`
	TotalCount int64           `json:"total_count"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	HasMore    bool            `json:"has_more"`
}
