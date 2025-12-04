package models

import "time"

// ============================================
// PROFILE MODELS
// ============================================

// Follow represents a follow relationship
type Follow struct {
	ID          string    `json:"id" db:"id"`
	FollowerID  string    `json:"follower_id" db:"follower_id"`
	FollowingID string    `json:"following_id" db:"following_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// PrivacySettings represents user privacy preferences
type PrivacySettings struct {
	UserID          string    `json:"user_id" db:"user_id"`
	AccountPrivacy  string    `json:"account_privacy" db:"account_privacy"`
	WhoCanComment   string    `json:"who_can_comment" db:"who_can_comment"`
	WhoCanDuet      string    `json:"who_can_duet" db:"who_can_duet"`
	WhoCanStitch    string    `json:"who_can_stitch" db:"who_can_stitch"`
	WhoCanMessage   string    `json:"who_can_message" db:"who_can_message"`
	ShowLikedVideos bool      `json:"show_liked_videos" db:"show_liked_videos"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// NotificationSettings represents user notification preferences
type NotificationSettings struct {
	UserID         string    `json:"user_id" db:"user_id"`
	Likes          bool      `json:"likes" db:"likes"`
	Comments       bool      `json:"comments" db:"comments"`
	NewFollowers   bool      `json:"new_followers" db:"new_followers"`
	Mentions       bool      `json:"mentions" db:"mentions"`
	LiveStreams    bool      `json:"live_streams" db:"live_streams"`
	DirectMessages bool      `json:"direct_messages" db:"direct_messages"`
	VideoUpdates   bool      `json:"video_updates" db:"video_updates"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// Referral represents a referral relationship
type Referral struct {
	ID           string    `json:"id" db:"id"`
	ReferrerID   string    `json:"referrer_id" db:"referrer_id"`
	ReferredID   string    `json:"referred_id" db:"referred_id"`
	ReferralCode string    `json:"referral_code" db:"referral_code"`
	RewardCoins  int       `json:"reward_coins" db:"reward_coins"`
	IsRewarded   bool      `json:"is_rewarded" db:"is_rewarded"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// ============================================
// REQUEST/RESPONSE DTOs
// ============================================

// EditProfileRequest represents a profile update request
type EditProfileRequest struct {
	Name             *string  `json:"name,omitempty"`
	FirstName        *string  `json:"first_name,omitempty"`
	LastName         *string  `json:"last_name,omitempty"`
	Username         *string  `json:"username,omitempty"`
	Bio              *string  `json:"bio,omitempty"`
	Website          *string  `json:"website,omitempty"`
	City             *string  `json:"city,omitempty"`
	Age              *int     `json:"age,omitempty"`
	Gender           *string  `json:"gender,omitempty"`
	RelationshipGoal *string  `json:"relationship_goal,omitempty"`
	Religion         *string  `json:"religion,omitempty"`
	Languages        []string `json:"languages,omitempty"`
	Interests        []string `json:"interests,omitempty"`
}

// FollowUserRequest represents a follow/unfollow request
type FollowUserRequest struct {
	UserID string `json:"user_id" validate:"required"`
	Action string `json:"action" validate:"required,oneof=follow unfollow"`
}

// UserProfileResponse represents a user profile with stats
type UserProfileResponse struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Age              int       `json:"age"`
	Gender           string    `json:"gender"`
	City             string    `json:"city"`
	Bio              *string   `json:"bio,omitempty"`
	Website          *string   `json:"website,omitempty"`
	IsVerified       bool      `json:"is_verified"`
	IsPrivate        bool      `json:"is_private"`
	FollowersCount   int       `json:"followers_count"`
	FollowingCount   int       `json:"following_count"`
	IsFollowing      bool      `json:"is_following"`
	IsFollowedBy     bool      `json:"is_followed_by"`
	IsBlocked        bool      `json:"is_blocked"`
	ReferralCode     *string   `json:"referral_code,omitempty"`
	RelationshipGoal string    `json:"relationship_goal"`
	Religion         *string   `json:"religion,omitempty"`
	Languages        []string  `json:"languages"`
	Interests        []string  `json:"interests"`
	CreatedAt        time.Time `json:"created_at"`
}

// FollowerResponse represents a follower/following user
type FollowerResponse struct {
	ID           string    `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Age          int       `json:"age" db:"age"`
	City         string    `json:"city" db:"city"`
	Bio          *string   `json:"bio,omitempty" db:"bio"`
	IsVerified   bool      `json:"is_verified" db:"is_verified"`
	IsFollowing  bool      `json:"is_following" db:"is_following"`
	IsFollowedBy bool      `json:"is_followed_by" db:"is_followed_by"`
	FollowedAt   time.Time `json:"followed_at" db:"followed_at"`
}

// UpdatePrivacySettingsRequest represents privacy settings update
type UpdatePrivacySettingsRequest struct {
	AccountPrivacy  *string `json:"account_privacy,omitempty" validate:"omitempty,oneof=public private"`
	WhoCanComment   *string `json:"who_can_comment,omitempty" validate:"omitempty,oneof=everyone followers nobody"`
	WhoCanDuet      *string `json:"who_can_duet,omitempty" validate:"omitempty,oneof=everyone followers nobody"`
	WhoCanStitch    *string `json:"who_can_stitch,omitempty" validate:"omitempty,oneof=everyone followers nobody"`
	WhoCanMessage   *string `json:"who_can_message,omitempty" validate:"omitempty,oneof=everyone followers nobody"`
	ShowLikedVideos *bool   `json:"show_liked_videos,omitempty"`
}

// UpdateNotificationSettingsRequest represents notification settings update
type UpdateNotificationSettingsRequest struct {
	Likes          *bool `json:"likes,omitempty"`
	Comments       *bool `json:"comments,omitempty"`
	NewFollowers   *bool `json:"new_followers,omitempty"`
	Mentions       *bool `json:"mentions,omitempty"`
	LiveStreams    *bool `json:"live_streams,omitempty"`
	DirectMessages *bool `json:"direct_messages,omitempty"`
	VideoUpdates   *bool `json:"video_updates,omitempty"`
}

// ApplyReferralCodeRequest represents referral code application
type ApplyReferralCodeRequest struct {
	ReferralCode string `json:"referral_code" validate:"required,min=8,max=20"`
}

// ReferralCodeResponse represents referral code info
type ReferralCodeResponse struct {
	ReferralCode   string `json:"referral_code" db:"referral_code"`
	TotalReferrals int64  `json:"total_referrals" db:"total_referrals"`
	TotalRewards   int64  `json:"total_rewards" db:"total_rewards"`
}

// BlockUserRequest represents a block/unblock request
type BlockUserRequest struct {
	UserID string `json:"user_id" validate:"required"`
	Action string `json:"action" validate:"required,oneof=block unblock"`
}

// BlockedUserResponse represents a blocked user
type BlockedUserResponse struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Age       int       `json:"age" db:"age"`
	City      string    `json:"city" db:"city"`
	BlockedAt time.Time `json:"blocked_at" db:"blocked_at"`
}

// VerificationRequest represents a verification request
type VerificationRequest struct {
	SelfieURL     string `json:"selfie_url" validate:"required,url"`
	IDDocumentURL string `json:"id_document_url" validate:"required,url"`
}

// ReportUserRequest represents a user report
type ReportUserRequest struct {
	UserID      string   `json:"user_id" validate:"required"`
	Reason      string   `json:"reason" validate:"required,oneof=inappropriate_content fake_profile harassment scam other"`
	Description string   `json:"description" validate:"required,min=10"`
	Screenshots []string `json:"screenshots,omitempty"`
}

// UserDetailResponse represents the full user details for the legacy app
type UserDetailResponse struct {
	ID              string                `json:"id" db:"id"`
	FirstName       string                `json:"first_name" db:"first_name"`
	LastName        string                `json:"last_name" db:"last_name"`
	Username        string                `json:"username" db:"username"`
	Email           string                `json:"email" db:"email"`
	Phone           string                `json:"phone" db:"phone"`
	Bio             string                `json:"bio" db:"bio"`
	Website         string                `json:"website" db:"website"`
	Gender          string                `json:"gender" db:"gender"`
	ProfilePic      string                `json:"profile_pic" db:"profile_pic"`
	ProfileGif      string                `json:"profile_gif" db:"profile_gif"`
	ProfileView     string                `json:"profile_view" db:"profile_view"`
	Wallet          int                   `json:"wallet" db:"wallet"`
	TotalBalanceUSD float64               `json:"total_balance_usd" db:"total_balance_usd"`
	TotalCoins      int                   `json:"total_coins_all_time" db:"total_coins_all_time"`
	FollowersCount  int                   `json:"followers_count" db:"followers_count"`
	FollowingCount  int                   `json:"following_count" db:"following_count"`
	LikesCount      int                   `json:"likes_count" db:"likes_count"`
	VideoCount      int                   `json:"video_count" db:"video_count"`
	Verified        int                   `json:"verified" db:"verified"`
	ReferralCode    string                `json:"referral_code" db:"referral_code"`
	PrivacySetting  *PrivacySettings      `json:"privacy_setting_model"`
	PushSetting     *NotificationSettings `json:"push_notification_model"`
}
