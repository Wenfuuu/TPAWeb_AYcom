package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Follow struct {
	ID 			uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	FollowerID	uuid.UUID `gorm:"type:uuid;not null" json:"follower_id"`
	FollowingID	uuid.UUID `gorm:"type:uuid;not null" json:"following_id"`
	CreatedAt	time.Time `json:"createdAt" gorm:"autoCreateTime"`
}

func (follow *Follow) BeforeCreate(tx *gorm.DB) (err error) {
	follow.ID = uuid.New()

	return
}