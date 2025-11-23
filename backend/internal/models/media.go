package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MediaType string

const (
	MediaTypePhoto MediaType = "photo"
	MediaTypeVideo MediaType = "video"
)

type Media struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	User      User      `gorm:"foreignKey:UserID"`

	MediaType     MediaType `gorm:"type:media_type;not null"`
	URL           string    `gorm:"type:text;not null"`
	ThumbnailURL  string    `gorm:"type:text"`
	DurationSeconds int     `gorm:"type:integer"`

	DisplayOrder int `gorm:"default:1;index"`

	IsApproved     bool   `gorm:"default:false"`
	ModerationNotes string `gorm:"type:text"`

	CreatedAt time.Time `gorm:"type:timestamptz;default:now()"`
	UpdatedAt time.Time `gorm:"type:timestamptz;default:now()"`
}

func (m *Media) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return
}

