package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReportReason string

const (
	ReportReasonInappropriateContent ReportReason = "inappropriate_content"
	ReportReasonFakeProfile          ReportReason = "fake_profile"
	ReportReasonHarassment           ReportReason = "harassment"
	ReportReasonScam                 ReportReason = "scam"
	ReportReasonOther                ReportReason = "other"
)

type Report struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ReporterID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Reporter      User      `gorm:"foreignKey:ReporterID"`
	ReportedUserID uuid.UUID `gorm:"type:uuid;not null;index"`
	ReportedUser  User      `gorm:"foreignKey:ReportedUserID"`

	Reason      ReportReason `gorm:"type:report_reason;not null"`
	Description string       `gorm:"type:text"`

	ScreenshotURLs JSONStringArray `gorm:"type:jsonb;default:'[]'"`

	IsReviewed bool       `gorm:"default:false;index"`
	ReviewedBy  *uuid.UUID `gorm:"type:uuid"`
	ReviewedAt  *time.Time `gorm:"type:timestamptz"`
	ActionTaken string     `gorm:"type:text"`

	CreatedAt time.Time `gorm:"type:timestamptz;default:now()"`
	UpdatedAt time.Time `gorm:"type:timestamptz;default:now()"`
}

func (r *Report) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return
}

