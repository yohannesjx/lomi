package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"lomi-backend/internal/models"

	"github.com/jmoiron/sqlx"
)

type ProfileRepository struct {
	db *sqlx.DB
}

func NewProfileRepository(db *sqlx.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

// ============================================
// PROFILE MANAGEMENT
// ============================================

// UpdateProfile updates user profile fields
func (r *ProfileRepository) UpdateProfile(ctx context.Context, userID string, req *models.EditProfileRequest) error {
	query := `UPDATE users SET `
	args := []interface{}{}
	argCount := 1
	updates := []string{}

	if req.Name != nil {
		updates = append(updates, fmt.Sprintf("name = $%d", argCount))
		args = append(args, *req.Name)
		argCount++
	}
	if req.Bio != nil {
		updates = append(updates, fmt.Sprintf("bio = $%d", argCount))
		args = append(args, *req.Bio)
		argCount++
	}
	if req.Website != nil {
		updates = append(updates, fmt.Sprintf("website = $%d", argCount))
		args = append(args, *req.Website)
		argCount++
	}
	if req.City != nil {
		updates = append(updates, fmt.Sprintf("city = $%d", argCount))
		args = append(args, *req.City)
		argCount++
	}
	if req.Age != nil {
		updates = append(updates, fmt.Sprintf("age = $%d", argCount))
		args = append(args, *req.Age)
		argCount++
	}
	if req.Gender != nil {
		updates = append(updates, fmt.Sprintf("gender = $%d::gender_type", argCount))
		args = append(args, *req.Gender)
		argCount++
	}
	if req.RelationshipGoal != nil {
		updates = append(updates, fmt.Sprintf("relationship_goal = $%d::relationship_goal", argCount))
		args = append(args, *req.RelationshipGoal)
		argCount++
	}
	if req.Religion != nil {
		updates = append(updates, fmt.Sprintf("religion = $%d::religion_type", argCount))
		args = append(args, *req.Religion)
		argCount++
	}
	if req.Languages != nil {
		languagesJSON, _ := json.Marshal(req.Languages)
		updates = append(updates, fmt.Sprintf("languages = $%d::jsonb", argCount))
		args = append(args, languagesJSON)
		argCount++
	}
	if req.Interests != nil {
		interestsJSON, _ := json.Marshal(req.Interests)
		updates = append(updates, fmt.Sprintf("interests = $%d::jsonb", argCount))
		args = append(args, interestsJSON)
		argCount++
	}

	if len(updates) == 0 {
		return nil // Nothing to update
	}

	query += fmt.Sprintf("%s WHERE id = $%d::uuid", joinStrings(updates, ", "), argCount)
	args = append(args, userID)

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

// ============================================
// FOLLOW MANAGEMENT
// ============================================

// FollowUser creates a follow relationship
func (r *ProfileRepository) FollowUser(ctx context.Context, followerID, followingID string) error {
	query := `
		INSERT INTO follows (follower_id, following_id)
		VALUES ($1::uuid, $2::uuid)
		ON CONFLICT (follower_id, following_id) DO NOTHING
	`
	_, err := r.db.ExecContext(ctx, query, followerID, followingID)
	return err
}

// UnfollowUser removes a follow relationship
func (r *ProfileRepository) UnfollowUser(ctx context.Context, followerID, followingID string) error {
	query := `DELETE FROM follows WHERE follower_id = $1::uuid AND following_id = $2::uuid`
	_, err := r.db.ExecContext(ctx, query, followerID, followingID)
	return err
}

// IsFollowing checks if user1 follows user2
func (r *ProfileRepository) IsFollowing(ctx context.Context, followerID, followingID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM follows WHERE follower_id = $1::uuid AND following_id = $2::uuid)`
	err := r.db.GetContext(ctx, &exists, query, followerID, followingID)
	return exists, err
}

// GetFollowers gets list of followers for a user
func (r *ProfileRepository) GetFollowers(ctx context.Context, userID string, viewerID string, page, pageSize int) ([]models.FollowerResponse, error) {
	offset := (page - 1) * pageSize
	query := `
		SELECT 
			u.id, u.name, u.age, u.city, u.bio, u.is_verified,
			EXISTS(SELECT 1 FROM follows WHERE follower_id = $2::uuid AND following_id = u.id) as is_following,
			EXISTS(SELECT 1 FROM follows WHERE follower_id = u.id AND following_id = $2::uuid) as is_followed_by,
			f.created_at as followed_at
		FROM follows f
		JOIN users u ON f.follower_id = u.id
		WHERE f.following_id = $1::uuid
		ORDER BY f.created_at DESC
		LIMIT $3 OFFSET $4
	`

	followers := make([]models.FollowerResponse, 0)
	err := r.db.SelectContext(ctx, &followers, query, userID, viewerID, pageSize, offset)
	if err != nil {
		return followers, err
	}
	return followers, nil
}

// GetFollowing gets list of users that a user follows
func (r *ProfileRepository) GetFollowing(ctx context.Context, userID string, viewerID string, page, pageSize int) ([]models.FollowerResponse, error) {
	offset := (page - 1) * pageSize
	query := `
		SELECT 
			u.id, u.name, u.age, u.city, u.bio, u.is_verified,
			EXISTS(SELECT 1 FROM follows WHERE follower_id = $2::uuid AND following_id = u.id) as is_following,
			EXISTS(SELECT 1 FROM follows WHERE follower_id = u.id AND following_id = $2::uuid) as is_followed_by,
			f.created_at as followed_at
		FROM follows f
		JOIN users u ON f.following_id = u.id
		WHERE f.follower_id = $1::uuid
		ORDER BY f.created_at DESC
		LIMIT $3 OFFSET $4
	`

	following := make([]models.FollowerResponse, 0)
	err := r.db.SelectContext(ctx, &following, query, userID, viewerID, pageSize, offset)
	if err != nil {
		return following, err
	}
	return following, nil
}

// ============================================
// BLOCK MANAGEMENT
// ============================================

// BlockUser creates a block relationship
func (r *ProfileRepository) BlockUser(ctx context.Context, blockerID, blockedID string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert block
	query := `
		INSERT INTO blocks (blocker_id, blocked_id)
		VALUES ($1::uuid, $2::uuid)
		ON CONFLICT (blocker_id, blocked_id) DO NOTHING
	`
	_, err = tx.ExecContext(ctx, query, blockerID, blockedID)
	if err != nil {
		return err
	}

	// Remove follow relationships
	_, err = tx.ExecContext(ctx, `DELETE FROM follows WHERE follower_id = $1::uuid AND following_id = $2::uuid`, blockerID, blockedID)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `DELETE FROM follows WHERE follower_id = $1::uuid AND following_id = $2::uuid`, blockedID, blockerID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// UnblockUser removes a block relationship
func (r *ProfileRepository) UnblockUser(ctx context.Context, blockerID, blockedID string) error {
	query := `DELETE FROM blocks WHERE blocker_id = $1::uuid AND blocked_id = $2::uuid`
	_, err := r.db.ExecContext(ctx, query, blockerID, blockedID)
	return err
}

// IsBlocked checks if user1 has blocked user2
func (r *ProfileRepository) IsBlocked(ctx context.Context, blockerID, blockedID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM blocks WHERE blocker_id = $1::uuid AND blocked_id = $2::uuid)`
	err := r.db.GetContext(ctx, &exists, query, blockerID, blockedID)
	return exists, err
}

// GetBlockedUsers gets list of blocked users
func (r *ProfileRepository) GetBlockedUsers(ctx context.Context, userID string, page, pageSize int) ([]models.BlockedUserResponse, error) {
	offset := (page - 1) * pageSize
	query := `
		SELECT 
			u.id, u.name, u.age, u.city,
			b.created_at as blocked_at
		FROM blocks b
		JOIN users u ON b.blocked_id = u.id
		WHERE b.blocker_id = $1::uuid
		ORDER BY b.created_at DESC
		LIMIT $2 OFFSET $3
	`

	blocked := make([]models.BlockedUserResponse, 0)
	err := r.db.SelectContext(ctx, &blocked, query, userID, pageSize, offset)
	if err != nil {
		return blocked, err
	}
	return blocked, nil
}

// ============================================
// PRIVACY SETTINGS
// ============================================

// GetPrivacySettings gets user's privacy settings
func (r *ProfileRepository) GetPrivacySettings(ctx context.Context, userID string) (*models.PrivacySettings, error) {
	var settings models.PrivacySettings
	query := `SELECT * FROM privacy_settings WHERE user_id = $1::uuid`
	err := r.db.GetContext(ctx, &settings, query, userID)
	if err == sql.ErrNoRows {
		// Create default settings
		return r.CreateDefaultPrivacySettings(ctx, userID)
	}
	return &settings, err
}

// CreateDefaultPrivacySettings creates default privacy settings for a user
func (r *ProfileRepository) CreateDefaultPrivacySettings(ctx context.Context, userID string) (*models.PrivacySettings, error) {
	query := `
		INSERT INTO privacy_settings (user_id)
		VALUES ($1::uuid)
		RETURNING *
	`
	var settings models.PrivacySettings
	err := r.db.GetContext(ctx, &settings, query, userID)
	return &settings, err
}

// UpdatePrivacySettings updates user's privacy settings
func (r *ProfileRepository) UpdatePrivacySettings(ctx context.Context, userID string, req *models.UpdatePrivacySettingsRequest) (*models.PrivacySettings, error) {
	query := `UPDATE privacy_settings SET `
	args := []interface{}{}
	argCount := 1
	updates := []string{}

	if req.AccountPrivacy != nil {
		updates = append(updates, fmt.Sprintf("account_privacy = $%d", argCount))
		args = append(args, *req.AccountPrivacy)
		argCount++
	}
	if req.WhoCanComment != nil {
		updates = append(updates, fmt.Sprintf("who_can_comment = $%d", argCount))
		args = append(args, *req.WhoCanComment)
		argCount++
	}
	if req.WhoCanDuet != nil {
		updates = append(updates, fmt.Sprintf("who_can_duet = $%d", argCount))
		args = append(args, *req.WhoCanDuet)
		argCount++
	}
	if req.WhoCanStitch != nil {
		updates = append(updates, fmt.Sprintf("who_can_stitch = $%d", argCount))
		args = append(args, *req.WhoCanStitch)
		argCount++
	}
	if req.WhoCanMessage != nil {
		updates = append(updates, fmt.Sprintf("who_can_message = $%d", argCount))
		args = append(args, *req.WhoCanMessage)
		argCount++
	}
	if req.ShowLikedVideos != nil {
		updates = append(updates, fmt.Sprintf("show_liked_videos = $%d", argCount))
		args = append(args, *req.ShowLikedVideos)
		argCount++
	}

	if len(updates) == 0 {
		return r.GetPrivacySettings(ctx, userID)
	}

	query += fmt.Sprintf("%s WHERE user_id = $%d::uuid RETURNING *", joinStrings(updates, ", "), argCount)
	args = append(args, userID)

	var settings models.PrivacySettings
	err := r.db.GetContext(ctx, &settings, query, args...)
	return &settings, err
}

// ============================================
// NOTIFICATION SETTINGS
// ============================================

// GetNotificationSettings gets user's notification settings
func (r *ProfileRepository) GetNotificationSettings(ctx context.Context, userID string) (*models.NotificationSettings, error) {
	var settings models.NotificationSettings
	query := `SELECT * FROM notification_settings WHERE user_id = $1::uuid`
	err := r.db.GetContext(ctx, &settings, query, userID)
	if err == sql.ErrNoRows {
		return r.CreateDefaultNotificationSettings(ctx, userID)
	}
	return &settings, err
}

// CreateDefaultNotificationSettings creates default notification settings
func (r *ProfileRepository) CreateDefaultNotificationSettings(ctx context.Context, userID string) (*models.NotificationSettings, error) {
	query := `
		INSERT INTO notification_settings (user_id)
		VALUES ($1::uuid)
		RETURNING *
	`
	var settings models.NotificationSettings
	err := r.db.GetContext(ctx, &settings, query, userID)
	return &settings, err
}

// UpdateNotificationSettings updates user's notification settings
func (r *ProfileRepository) UpdateNotificationSettings(ctx context.Context, userID string, req *models.UpdateNotificationSettingsRequest) (*models.NotificationSettings, error) {
	query := `UPDATE notification_settings SET `
	args := []interface{}{}
	argCount := 1
	updates := []string{}

	if req.Likes != nil {
		updates = append(updates, fmt.Sprintf("likes = $%d", argCount))
		args = append(args, *req.Likes)
		argCount++
	}
	if req.Comments != nil {
		updates = append(updates, fmt.Sprintf("comments = $%d", argCount))
		args = append(args, *req.Comments)
		argCount++
	}
	if req.NewFollowers != nil {
		updates = append(updates, fmt.Sprintf("new_followers = $%d", argCount))
		args = append(args, *req.NewFollowers)
		argCount++
	}
	if req.Mentions != nil {
		updates = append(updates, fmt.Sprintf("mentions = $%d", argCount))
		args = append(args, *req.Mentions)
		argCount++
	}
	if req.LiveStreams != nil {
		updates = append(updates, fmt.Sprintf("live_streams = $%d", argCount))
		args = append(args, *req.LiveStreams)
		argCount++
	}
	if req.DirectMessages != nil {
		updates = append(updates, fmt.Sprintf("direct_messages = $%d", argCount))
		args = append(args, *req.DirectMessages)
		argCount++
	}
	if req.VideoUpdates != nil {
		updates = append(updates, fmt.Sprintf("video_updates = $%d", argCount))
		args = append(args, *req.VideoUpdates)
		argCount++
	}

	if len(updates) == 0 {
		return r.GetNotificationSettings(ctx, userID)
	}

	query += fmt.Sprintf("%s WHERE user_id = $%d::uuid RETURNING *", joinStrings(updates, ", "), argCount)
	args = append(args, userID)

	var settings models.NotificationSettings
	err := r.db.GetContext(ctx, &settings, query, args...)
	return &settings, err
}

// ============================================
// REFERRAL SYSTEM
// ============================================

// GetReferralCode gets user's referral code
func (r *ProfileRepository) GetReferralCode(ctx context.Context, userID string) (string, error) {
	var code string
	query := `SELECT referral_code FROM users WHERE id = $1::uuid`
	err := r.db.GetContext(ctx, &code, query, userID)
	return code, err
}

// GetReferralStats gets referral statistics
func (r *ProfileRepository) GetReferralStats(ctx context.Context, userID string) (*models.ReferralCodeResponse, error) {
	var stats models.ReferralCodeResponse
	query := `
		SELECT 
			u.referral_code,
			COUNT(r.id) as total_referrals,
			COALESCE(SUM(r.reward_coins), 0) as total_rewards
		FROM users u
		LEFT JOIN referrals r ON u.id = r.referrer_id
		WHERE u.id = $1::uuid
		GROUP BY u.referral_code
	`
	err := r.db.GetContext(ctx, &stats, query, userID)
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

// ApplyReferralCode applies a referral code for a new user
func (r *ProfileRepository) ApplyReferralCode(ctx context.Context, userID, referralCode string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get referrer ID
	var referrerID string
	err = tx.GetContext(ctx, &referrerID, `SELECT id FROM users WHERE referral_code = $1`, referralCode)
	if err != nil {
		return fmt.Errorf("invalid referral code")
	}

	// Check if user already used a referral code
	var exists bool
	err = tx.GetContext(ctx, &exists, `SELECT EXISTS(SELECT 1 FROM referrals WHERE referred_id = $1::uuid)`, userID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("referral code already applied")
	}

	// Create referral record
	rewardCoins := 50 // Default reward
	query := `
		INSERT INTO referrals (referrer_id, referred_id, referral_code, reward_coins, is_rewarded)
		VALUES ($1::uuid, $2::uuid, $3, $4, true)
	`
	_, err = tx.ExecContext(ctx, query, referrerID, userID, referralCode, rewardCoins)
	if err != nil {
		return err
	}

	// Award coins to referrer
	_, err = tx.ExecContext(ctx, `UPDATE users SET coin_balance = coin_balance + $1 WHERE id = $2::uuid`, rewardCoins, referrerID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// ============================================
// ACCOUNT MANAGEMENT
// ============================================

// DeleteUserAccount soft deletes a user account
func (r *ProfileRepository) DeleteUserAccount(ctx context.Context, userID string) error {
	query := `UPDATE users SET deleted_at = NOW() WHERE id = $1::uuid AND deleted_at IS NULL`
	result, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found or already deleted")
	}

	return nil
}

// RequestVerification creates a verification request
func (r *ProfileRepository) RequestVerification(ctx context.Context, userID, selfieURL, idDocumentURL string) error {
	// Check if user already has a pending or approved verification
	var exists bool
	checkQuery := `
		SELECT EXISTS(
			SELECT 1 FROM verifications 
			WHERE user_id = $1::uuid 
			AND status IN ('pending', 'approved')
		)
	`
	err := r.db.GetContext(ctx, &exists, checkQuery, userID)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("verification already requested or approved")
	}

	// Create verification request
	query := `
		INSERT INTO verifications (user_id, selfie_url, id_document_url, status)
		VALUES ($1::uuid, $2, $3, 'pending')
	`
	_, err = r.db.ExecContext(ctx, query, userID, selfieURL, idDocumentURL)
	return err
}

// ReportUser creates a report against a user
func (r *ProfileRepository) ReportUser(ctx context.Context, reporterID, reportedUserID, reason, description string, screenshots []string) error {
	// Check if user already reported this user recently (within 24 hours)
	var exists bool
	checkQuery := `
		SELECT EXISTS(
			SELECT 1 FROM reports 
			WHERE reporter_id = $1::uuid 
			AND reported_user_id = $2::uuid
			AND created_at > NOW() - INTERVAL '24 hours'
		)
	`
	err := r.db.GetContext(ctx, &exists, checkQuery, reporterID, reportedUserID)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("you have already reported this user recently")
	}

	// Convert screenshots to JSONB
	screenshotsJSON := "[]"
	if len(screenshots) > 0 {
		bytes, _ := json.Marshal(screenshots)
		screenshotsJSON = string(bytes)
	}

	// Create report
	query := `
		INSERT INTO reports (reporter_id, reported_user_id, reason, description, screenshot_urls)
		VALUES ($1::uuid, $2::uuid, $3::report_reason, $4, $5::jsonb)
	`
	_, err = r.db.ExecContext(ctx, query, reporterID, reportedUserID, reason, description, screenshotsJSON)
	return err
}

// ============================================
// SOCIAL FEATURES
// ============================================

// TrackProfileShare tracks when a user shares a profile
func (r *ProfileRepository) TrackProfileShare(ctx context.Context, userID, sharedBy, platform string) error {
	query := `
		INSERT INTO profile_shares (user_id, shared_by, platform)
		VALUES ($1::uuid, $2::uuid, $3)
	`
	_, err := r.db.ExecContext(ctx, query, userID, sharedBy, platform)
	return err
}

// GetProfileShareCount gets total shares for a user
func (r *ProfileRepository) GetProfileShareCount(ctx context.Context, userID string) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM profile_shares WHERE user_id = $1::uuid`
	err := r.db.GetContext(ctx, &count, query, userID)
	return count, err
}

// ============================================
// HELPER FUNCTIONS
// ============================================

func joinStrings(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}

// ============================================
// APP SETTINGS
// ============================================

// ChangeAppLanguage changes user's app language
func (r *ProfileRepository) ChangeAppLanguage(ctx context.Context, userID, language string) error {
	query := `UPDATE users SET app_language = $1 WHERE id = $2::uuid`
	_, err := r.db.ExecContext(ctx, query, language, userID)
	return err
}

// ChangeAppTheme changes user's app theme
func (r *ProfileRepository) ChangeAppTheme(ctx context.Context, userID, theme string) error {
	query := `UPDATE users SET app_theme = $1 WHERE id = $2::uuid`
	_, err := r.db.ExecContext(ctx, query, theme, userID)
	return err
}

// ClearCache marks cache as cleared
func (r *ProfileRepository) ClearCache(ctx context.Context, userID string) error {
	query := `UPDATE users SET cache_cleared_at = NOW() WHERE id = $1::uuid`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

// GetAppSettings gets user's app settings
func (r *ProfileRepository) GetAppSettings(ctx context.Context, userID string) (*models.AppSettings, error) {
	var settings models.AppSettings
	query := `SELECT app_language, app_theme, cache_cleared_at FROM users WHERE id = $1::uuid`
	err := r.db.GetContext(ctx, &settings, query, userID)
	return &settings, err
}
