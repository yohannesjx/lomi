package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageType string

const (
	MessageTypeText       MessageType = "text"
	MessageTypePhoto      MessageType = "photo"
	MessageTypeVideo      MessageType = "video"
	MessageTypeVoice      MessageType = "voice"
	MessageTypeSticker    MessageType = "sticker"
	MessageTypeGift       MessageType = "gift"
	MessageTypeBunaInvite MessageType = "buna_invite"
)

type Message struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`

	// Private chat fields
	MatchID    *uuid.UUID `gorm:"type:uuid;index"`
	Match      *Match     `gorm:"foreignKey:MatchID"`
	ReceiverID *uuid.UUID `gorm:"type:uuid;index"`
	Receiver   *User      `gorm:"foreignKey:ReceiverID"`

	// Live chat fields
	LiveStreamID *uuid.UUID `gorm:"type:uuid;index"`
	LiveStream   *User      `gorm:"foreignKey:LiveStreamID"`
	IsLive       bool       `gorm:"default:false;not null;index"`
	IsSystem     bool       `gorm:"default:false;not null"`
	Seq          int64      `gorm:"default:0"`
	Pinned       bool       `gorm:"default:false;not null"`

	// Common fields
	SenderID uuid.UUID `gorm:"type:uuid;not null;index"`
	Sender   User      `gorm:"foreignKey:SenderID"`

	MessageType MessageType `gorm:"type:message_type;not null;default:'text'"`

	Content  string     `gorm:"type:text"`
	MediaURL string     `gorm:"type:text"`
	GiftID   *uuid.UUID `gorm:"type:uuid"`
	Gift     *Gift      `gorm:"foreignKey:GiftID"`

	Metadata JSONMap    `gorm:"type:jsonb;default:'{}'"`
	IsRead   bool       `gorm:"default:false;index"`
	ReadAt   *time.Time `gorm:"type:timestamptz"`

	CreatedAt time.Time `gorm:"type:timestamptz;default:now()"`
	UpdatedAt time.Time `gorm:"type:timestamptz;default:now()"`
}

func (m *Message) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return
}
