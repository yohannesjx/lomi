package services

import (
	"context"
	"fmt"

	"lomi-backend/internal/models"
	"lomi-backend/internal/repositories"
)

type ProfileService struct {
	profileRepo *repositories.ProfileRepository
}

func NewProfileService(profileRepo *repositories.ProfileRepository) *ProfileService {
	return &ProfileService{
		profileRepo: profileRepo,
	}
}

// ============================================
// PROFILE MANAGEMENT
// ============================================

// UpdateProfile updates user profile
func (s *ProfileService) UpdateProfile(ctx context.Context, userID string, req *models.EditProfileRequest) error {
	// Validate age if provided
	if req.Age != nil && (*req.Age < 18 || *req.Age > 100) {
		return fmt.Errorf("age must be between 18 and 100")
	}

	// Validate gender if provided
	if req.Gender != nil {
		validGenders := map[string]bool{"male": true, "female": true, "other": true}
		if !validGenders[*req.Gender] {
			return fmt.Errorf("invalid gender")
		}
	}

	// Validate relationship goal if provided
	if req.RelationshipGoal != nil {
		validGoals := map[string]bool{"friends": true, "dating": true, "serious": true}
		if !validGoals[*req.RelationshipGoal] {
			return fmt.Errorf("invalid relationship goal")
		}
	}

	// Validate religion if provided
	if req.Religion != nil {
		validReligions := map[string]bool{
			"orthodox": true, "muslim": true, "protestant": true,
			"catholic": true, "other": true, "none": true,
		}
		if !validReligions[*req.Religion] {
			return fmt.Errorf("invalid religion")
		}
	}

	return s.profileRepo.UpdateProfile(ctx, userID, req)
}

// ============================================
// FOLLOW MANAGEMENT
// ============================================

// FollowUser follows or unfollows a user
func (s *ProfileService) FollowUser(ctx context.Context, followerID string, req *models.FollowUserRequest) error {
	// Validate not following self
	if followerID == req.UserID {
		return fmt.Errorf("cannot follow yourself")
	}

	// Check if blocked
	blocked, err := s.profileRepo.IsBlocked(ctx, req.UserID, followerID)
	if err != nil {
		return err
	}
	if blocked {
		return fmt.Errorf("you are blocked by this user")
	}

	blocked, err = s.profileRepo.IsBlocked(ctx, followerID, req.UserID)
	if err != nil {
		return err
	}
	if blocked {
		return fmt.Errorf("you have blocked this user")
	}

	if req.Action == "follow" {
		return s.profileRepo.FollowUser(ctx, followerID, req.UserID)
	} else if req.Action == "unfollow" {
		return s.profileRepo.UnfollowUser(ctx, followerID, req.UserID)
	}

	return fmt.Errorf("invalid action")
}

// GetFollowers gets list of followers
func (s *ProfileService) GetFollowers(ctx context.Context, userID, viewerID string, page, pageSize int) ([]models.FollowerResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	return s.profileRepo.GetFollowers(ctx, userID, viewerID, page, pageSize)
}

// GetFollowing gets list of following
func (s *ProfileService) GetFollowing(ctx context.Context, userID, viewerID string, page, pageSize int) ([]models.FollowerResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	return s.profileRepo.GetFollowing(ctx, userID, viewerID, page, pageSize)
}

// ============================================
// BLOCK MANAGEMENT
// ============================================

// BlockUser blocks or unblocks a user
func (s *ProfileService) BlockUser(ctx context.Context, blockerID string, req *models.BlockUserRequest) error {
	// Validate not blocking self
	if blockerID == req.UserID {
		return fmt.Errorf("cannot block yourself")
	}

	if req.Action == "block" {
		return s.profileRepo.BlockUser(ctx, blockerID, req.UserID)
	} else if req.Action == "unblock" {
		return s.profileRepo.UnblockUser(ctx, blockerID, req.UserID)
	}

	return fmt.Errorf("invalid action")
}

// GetBlockedUsers gets list of blocked users
func (s *ProfileService) GetBlockedUsers(ctx context.Context, userID string, page, pageSize int) ([]models.BlockedUserResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	return s.profileRepo.GetBlockedUsers(ctx, userID, page, pageSize)
}

// ============================================
// PRIVACY SETTINGS
// ============================================

// GetPrivacySettings gets user's privacy settings
func (s *ProfileService) GetPrivacySettings(ctx context.Context, userID string) (*models.PrivacySettings, error) {
	return s.profileRepo.GetPrivacySettings(ctx, userID)
}

// UpdatePrivacySettings updates user's privacy settings
func (s *ProfileService) UpdatePrivacySettings(ctx context.Context, userID string, req *models.UpdatePrivacySettingsRequest) (*models.PrivacySettings, error) {
	return s.profileRepo.UpdatePrivacySettings(ctx, userID, req)
}

// ============================================
// NOTIFICATION SETTINGS
// ============================================

// GetNotificationSettings gets user's notification settings
func (s *ProfileService) GetNotificationSettings(ctx context.Context, userID string) (*models.NotificationSettings, error) {
	return s.profileRepo.GetNotificationSettings(ctx, userID)
}

// UpdateNotificationSettings updates user's notification settings
func (s *ProfileService) UpdateNotificationSettings(ctx context.Context, userID string, req *models.UpdateNotificationSettingsRequest) (*models.NotificationSettings, error) {
	return s.profileRepo.UpdateNotificationSettings(ctx, userID, req)
}

// ============================================
// REFERRAL SYSTEM
// ============================================

// GetReferralCode gets user's referral code and stats
func (s *ProfileService) GetReferralCode(ctx context.Context, userID string) (*models.ReferralCodeResponse, error) {
	return s.profileRepo.GetReferralStats(ctx, userID)
}

// ApplyReferralCode applies a referral code
func (s *ProfileService) ApplyReferralCode(ctx context.Context, userID string, req *models.ApplyReferralCodeRequest) error {
	if len(req.ReferralCode) < 8 || len(req.ReferralCode) > 20 {
		return fmt.Errorf("invalid referral code format")
	}

	return s.profileRepo.ApplyReferralCode(ctx, userID, req.ReferralCode)
}

// ============================================
// ACCOUNT MANAGEMENT
// ============================================

// DeleteUserAccount deletes a user account
func (s *ProfileService) DeleteUserAccount(ctx context.Context, userID string) error {
	return s.profileRepo.DeleteUserAccount(ctx, userID)
}

// RequestVerification requests verification badge
func (s *ProfileService) RequestVerification(ctx context.Context, userID, selfieURL, idDocumentURL string) error {
	if selfieURL == "" || idDocumentURL == "" {
		return fmt.Errorf("selfie and ID document are required")
	}

	return s.profileRepo.RequestVerification(ctx, userID, selfieURL, idDocumentURL)
}

// ReportUser reports a user
func (s *ProfileService) ReportUser(ctx context.Context, reporterID, reportedUserID, reason, description string, screenshots []string) error {
	// Validate not reporting self
	if reporterID == reportedUserID {
		return fmt.Errorf("cannot report yourself")
	}

	// Validate reason
	validReasons := map[string]bool{
		"inappropriate_content": true,
		"fake_profile":          true,
		"harassment":            true,
		"scam":                  true,
		"other":                 true,
	}
	if !validReasons[reason] {
		return fmt.Errorf("invalid report reason")
	}

	if description == "" {
		return fmt.Errorf("description is required")
	}

	return s.profileRepo.ReportUser(ctx, reporterID, reportedUserID, reason, description, screenshots)
}

// ============================================
// SOCIAL FEATURES
// ============================================

// GenerateQRCode generates QR code data for a user profile
func (s *ProfileService) GenerateQRCode(ctx context.Context, userID string) (*models.QRCodeResponse, error) {
	// Get user's referral code
	code, err := s.profileRepo.GetReferralCode(ctx, userID)
	if err != nil {
		return nil, err
	}

	// In production, you would generate actual QR code image here
	// For now, return the data needed to generate QR code on client side
	profileURL := fmt.Sprintf("https://lomi.app/profile/%s", userID)
	qrCodeURL := fmt.Sprintf("https://api.qrserver.com/v1/create-qr-code/?size=300x300&data=%s", profileURL)

	return &models.QRCodeResponse{
		QRCodeURL:    qrCodeURL,
		ProfileURL:   profileURL,
		ReferralCode: code,
	}, nil
}

// ShareProfile tracks a profile share
func (s *ProfileService) ShareProfile(ctx context.Context, sharedBy string, req *models.ShareProfileRequest) error {
	// Validate platform
	validPlatforms := map[string]bool{
		"whatsapp": true, "telegram": true, "twitter": true,
		"facebook": true, "instagram": true, "link": true,
		"sms": true, "email": true, "other": true,
	}
	if !validPlatforms[req.Platform] {
		return fmt.Errorf("invalid platform")
	}

	return s.profileRepo.TrackProfileShare(ctx, req.UserID, sharedBy, req.Platform)
}

// ============================================
// APP SETTINGS
// ============================================

// ChangeAppLanguage changes user's app language
func (s *ProfileService) ChangeAppLanguage(ctx context.Context, userID, language string) error {
	validLanguages := map[string]bool{
		"en": true, "am": true, "om": true, "ti": true, "so": true, "ar": true,
	}
	if !validLanguages[language] {
		return fmt.Errorf("invalid language")
	}
	return s.profileRepo.ChangeAppLanguage(ctx, userID, language)
}

// ChangeAppTheme changes user's app theme
func (s *ProfileService) ChangeAppTheme(ctx context.Context, userID, theme string) error {
	validThemes := map[string]bool{"light": true, "dark": true, "auto": true}
	if !validThemes[theme] {
		return fmt.Errorf("invalid theme")
	}
	return s.profileRepo.ChangeAppTheme(ctx, userID, theme)
}

// ClearCache clears user's cache
func (s *ProfileService) ClearCache(ctx context.Context, userID string) error {
	return s.profileRepo.ClearCache(ctx, userID)
}

// GetAppSettings gets user's app settings
func (s *ProfileService) GetAppSettings(ctx context.Context, userID string) (*models.AppSettings, error) {
	return s.profileRepo.GetAppSettings(ctx, userID)
}
