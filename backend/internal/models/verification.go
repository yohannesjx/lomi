package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Verification struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID uuid.UUID `gorm:"type:uuid;not null;index"`
	User   User      `gorm:"foreignKey:UserID"`

	SelfieURL       string `gorm:"type:text;not null"`
	IDDocumentURL   string `gorm:"type:text;not null"`

	Status        VerificationStatus `gorm:"type:verification_status;default:'pending';index"`
	ReviewedBy    *uuid.UUID          `gorm:"type:uuid"`
	ReviewedAt    *time.Time          `gorm:"type:timestamptz"`
	RejectionReason string            `gorm:"type:text"`

	CreatedAt time.Time `gorm:"type:timestamptz;default:now()"`
	UpdatedAt time.Time `gorm:"type:timestamptz;default:now()"`
}

func (v *Verification) BeforeCreate(tx *gorm.DB) (err error) {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	return
}

