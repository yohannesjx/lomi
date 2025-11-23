package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PayoutStatus string

const (
	PayoutStatusPending    PayoutStatus = "pending"
	PayoutStatusProcessing PayoutStatus = "processing"
	PayoutStatusCompleted  PayoutStatus = "completed"
	PayoutStatusRejected   PayoutStatus = "rejected"
)

type Payout struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID uuid.UUID `gorm:"type:uuid;not null;index"`
	User   User      `gorm:"foreignKey:UserID"`

	GiftBalanceAmount    float64 `gorm:"type:decimal(10,2);not null"`
	PlatformFeePercentage int    `gorm:"not null;default:25"`
	PlatformFeeAmount    float64 `gorm:"type:decimal(10,2);not null"`
	NetAmount            float64 `gorm:"type:decimal(10,2);not null"`

	PaymentMethod      PaymentMethod `gorm:"type:payment_method;not null"`
	PaymentAccount     string        `gorm:"size:255;not null"`
	PaymentAccountName string        `gorm:"size:255"`

	Status        PayoutStatus `gorm:"type:payout_status;default:'pending';index"`
	ProcessedBy   *uuid.UUID   `gorm:"type:uuid"`
	ProcessedAt   *time.Time   `gorm:"type:timestamptz"`
	PaymentReference string     `gorm:"size:255"`

	AdminNotes      string `gorm:"type:text"`
	RejectionReason string `gorm:"type:text"`

	CreatedAt time.Time `gorm:"type:timestamptz;default:now()"`
	UpdatedAt time.Time `gorm:"type:timestamptz;default:now()"`
}

func (p *Payout) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}

