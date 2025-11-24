package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionType string
type PaymentMethod string
type PaymentStatus string

const (
	TransactionTypePurchase                  TransactionType = "purchase"
	TransactionTypeGiftSent                  TransactionType = "gift_sent"
	TransactionTypeGiftReceived              TransactionType = "gift_received"
	TransactionTypeBoost                     TransactionType = "boost"
	TransactionTypeRefund                    TransactionType = "refund"
	TransactionTypeChannelSubscriptionReward TransactionType = "channel_subscription_reward"
	TransactionTypeReveal                    TransactionType = "reveal"

	PaymentMethodTelebirr  PaymentMethod = "telebirr"
	PaymentMethodCbeBirr   PaymentMethod = "cbe_birr"
	PaymentMethodHelloCash PaymentMethod = "hellocash"
	PaymentMethodAmole     PaymentMethod = "amole"

	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

type CoinTransaction struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID uuid.UUID `gorm:"type:uuid;not null;index"`
	User   User      `gorm:"foreignKey:UserID"`

	TransactionType TransactionType `gorm:"type:transaction_type;not null;index"`

	CoinAmount int `gorm:"not null"`

	// For purchases
	BirrAmount       float64       `gorm:"type:decimal(10,2)"`
	PaymentMethod    PaymentMethod `gorm:"type:payment_method"`
	PaymentReference string        `gorm:"size:255"`
	PaymentStatus    PaymentStatus `gorm:"type:payment_status;default:'pending';index"`

	// For gifts
	GiftTransactionID *uuid.UUID       `gorm:"type:uuid"`
	GiftTransaction   *GiftTransaction `gorm:"foreignKey:GiftTransactionID"`

	BalanceAfter int `gorm:"not null"`

	Metadata JSONMap `gorm:"type:jsonb;default:'{}'"`

	CreatedAt time.Time `gorm:"type:timestamptz;default:now()"`
	UpdatedAt time.Time `gorm:"type:timestamptz;default:now()"`
}

func (ct *CoinTransaction) BeforeCreate(tx *gorm.DB) (err error) {
	if ct.ID == uuid.Nil {
		ct.ID = uuid.New()
	}
	return
}

type GiftTransaction struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	SenderID   uuid.UUID `gorm:"type:uuid;not null;index"`
	Sender     User      `gorm:"foreignKey:SenderID"`
	ReceiverID uuid.UUID `gorm:"type:uuid;not null;index"`
	Receiver   User      `gorm:"foreignKey:ReceiverID"`
	GiftID     uuid.UUID `gorm:"type:uuid;not null;index"`
	Gift       Gift      `gorm:"foreignKey:GiftID"`

	CoinAmount int     `gorm:"not null"`
	BirrValue  float64 `gorm:"type:decimal(10,2);not null"`
	GiftType   string  `gorm:"size:50"` // e.g., "rose", "universe", "lomi_crown"

	MessageID *uuid.UUID `gorm:"type:uuid"`
	Message   *Message   `gorm:"foreignKey:MessageID"`

	CreatedAt time.Time `gorm:"type:timestamptz;default:now()"`
}

func (gt *GiftTransaction) BeforeCreate(tx *gorm.DB) (err error) {
	if gt.ID == uuid.Nil {
		gt.ID = uuid.New()
	}
	return
}
