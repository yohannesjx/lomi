package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Gift struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	NameEn      string    `gorm:"size:255;not null"`
	NameAm      string    `gorm:"size:255;not null"`
	DescriptionEn string  `gorm:"type:text"`
	DescriptionAm string  `gorm:"type:text"`

	CoinPrice int     `gorm:"not null;check:coin_price > 0"`
	BirrValue float64 `gorm:"type:decimal(10,2);not null;check:birr_value > 0"`

	IconURL      string `gorm:"type:text;not null"`
	AnimationURL string `gorm:"type:text;not null"`
	SoundURL     string `gorm:"type:text"`

	HasSpecialEffect        bool `gorm:"default:false"`
	SpecialEffectDurationDays int `gorm:"type:integer"`

	IsActive  bool `gorm:"default:true;index"`
	IsFeatured bool `gorm:"default:false;index"`
	DisplayOrder int `gorm:"default:0;index"`

	CreatedAt time.Time `gorm:"type:timestamptz;default:now()"`
	UpdatedAt time.Time `gorm:"type:timestamptz;default:now()"`
}

func (g *Gift) BeforeCreate(tx *gorm.DB) (err error) {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	return
}

