package models

import "time"

// ============================================
// SOCIAL FEATURES MODELS
// ============================================

// ProfileShare represents a profile share event
type ProfileShare struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	SharedBy  string    `json:"shared_by" db:"shared_by"`
	Platform  string    `json:"platform" db:"platform"` // whatsapp, telegram, twitter, etc.
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// QRCodeResponse represents QR code data
type QRCodeResponse struct {
	QRCodeURL    string `json:"qr_code_url"`
	ProfileURL   string `json:"profile_url"`
	ReferralCode string `json:"referral_code,omitempty"`
}

// ShareProfileRequest represents a share tracking request
type ShareProfileRequest struct {
	UserID   string `json:"user_id" validate:"required"`
	Platform string `json:"platform" validate:"required,oneof=whatsapp telegram twitter facebook instagram link sms email other"`
}
