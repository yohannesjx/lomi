package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Match struct {
	ID      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	User1ID uuid.UUID `gorm:"type:uuid;not null;index"`
	User1   User      `gorm:"foreignKey:User1ID"`
	User2ID uuid.UUID `gorm:"type:uuid;not null;index"`
	User2   User      `gorm:"foreignKey:User2ID"`

	InitiatedBy uuid.UUID `gorm:"type:uuid;not null"`
	Initiator   User      `gorm:"foreignKey:InitiatedBy"`

	IsActive   bool       `gorm:"default:true;index"`
	UnmatchedBy *uuid.UUID `gorm:"type:uuid"`
	UnmatchedAt *time.Time `gorm:"type:timestamptz"`

	CreatedAt time.Time `gorm:"type:timestamptz;default:now()"`
}

func (m *Match) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	// Ensure consistent ordering (user1_id < user2_id)
	if m.User1ID.String() > m.User2ID.String() {
		m.User1ID, m.User2ID = m.User2ID, m.User1ID
	}
	return
}

