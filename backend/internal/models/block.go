package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Block struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	BlockerID uuid.UUID `gorm:"type:uuid;not null;index"`
	Blocker   User      `gorm:"foreignKey:BlockerID"`
	BlockedID uuid.UUID `gorm:"type:uuid;not null;index"`
	Blocked   User      `gorm:"foreignKey:BlockedID"`

	CreatedAt time.Time `gorm:"type:timestamptz;default:now()"`
}

func (b *Block) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return
}

