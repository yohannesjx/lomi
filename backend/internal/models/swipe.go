package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SwipeAction string

const (
	SwipeActionLike     SwipeAction = "like"
	SwipeActionPass     SwipeAction = "pass"
	SwipeActionSuperLike SwipeAction = "super_like"
)

type Swipe struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	SwiperID uuid.UUID `gorm:"type:uuid;not null;index"`
	Swiper   User      `gorm:"foreignKey:SwiperID"`
	SwipedID uuid.UUID `gorm:"type:uuid;not null;index"`
	Swiped   User      `gorm:"foreignKey:SwipedID"`

	Action SwipeAction `gorm:"type:swipe_action;not null"`

	CreatedAt time.Time `gorm:"type:timestamptz;default:now()"`
}

func (s *Swipe) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}

