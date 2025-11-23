package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RewardChannel struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ChannelUsername string  `gorm:"size:255;not null"`
	ChannelName   string    `gorm:"size:255;not null"`
	ChannelLink   string    `gorm:"size:255;not null"`

	CoinReward int `gorm:"not null;default:50"`

	IconURL string `gorm:"type:text"`

	IsActive    bool `gorm:"default:true;index"`
	DisplayOrder int `gorm:"default:0;index"`

	CreatedAt time.Time `gorm:"type:timestamptz;default:now()"`
	UpdatedAt time.Time `gorm:"type:timestamptz;default:now()"`
}

func (rc *RewardChannel) BeforeCreate(tx *gorm.DB) (err error) {
	if rc.ID == uuid.Nil {
		rc.ID = uuid.New()
	}
	return
}

type UserChannelReward struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;index"`
	User       User      `gorm:"foreignKey:UserID"`
	ChannelID  uuid.UUID `gorm:"type:uuid;not null;index"`
	Channel    RewardChannel `gorm:"foreignKey:ChannelID"`

	RewardAmount int `gorm:"not null"`

	ClaimedAt time.Time `gorm:"type:timestamptz;default:now()"`
}

func (ucr *UserChannelReward) BeforeCreate(tx *gorm.DB) (err error) {
	if ucr.ID == uuid.Nil {
		ucr.ID = uuid.New()
	}
	return
}

