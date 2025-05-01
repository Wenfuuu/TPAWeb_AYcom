package database

import (
	"github.com/Wenfuuu/tpa-web/service-follow/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	followerId, err := uuid.Parse("ae2cf219-a6e3-46d7-aee0-6f05eb23f602")
	if err != nil {
		return
	}

	followingId, err := uuid.Parse("f3122b31-7b92-400d-b94d-28ef100ca23b")
	if err != nil {
		return
	}

	var follows = []model.Follow{
		{
			FollowerID:     followerId,
			FollowingID:	followingId,
		},
	}

	for _, element := range follows {
		db.Create(&element)
	}
}