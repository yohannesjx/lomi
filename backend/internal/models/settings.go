package models

import (
	"time"

	"github.com/google/uuid"
)

// PrivacySetting Model
type PrivacySetting struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID         uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	VideosDownload int       `gorm:"default:1"`
	DirectMessage  int       `gorm:"default:1"`
	Duet           int       `gorm:"default:1"`
	LikedVideos    int       `gorm:"default:1"`
	VideoComment   int       `gorm:"default:1"`
	OrderHistory   int       `gorm:"default:1"`
	CreatedAt      time.Time `gorm:"type:timestamptz;default:now()"`
	UpdatedAt      time.Time `gorm:"type:timestamptz;default:now()"`
}

// PushNotification Model
type PushNotification struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID         uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	Likes          int       `gorm:"default:1"`
	Comments       int       `gorm:"default:1"`
	NewFollowers   int       `gorm:"default:1"`
	Mentions       int       `gorm:"default:1"`
	DirectMessages int       `gorm:"default:1"`
	VideoUpdates   int       `gorm:"default:1"`
	CreatedAt      time.Time `gorm:"type:timestamptz;default:now()"`
	UpdatedAt      time.Time `gorm:"type:timestamptz;default:now()"`
}
