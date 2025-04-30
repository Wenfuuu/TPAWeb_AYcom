package database

import (
	"time"

	"github.com/Wenfuuu/tpa-web/service-thread/model"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	var threads = []model.Thread{
		{
			UserID:          uuid.New(),
			CategoryID:      uuid.New(),
			CommunityID:     uuid.New(),
			Content:         "This is my first thread. Excited to be here!",
			ReplyPermission: "everyone",
			IsAdvertisement: false,
			Status:          "published",
			CreatedAt:       time.Now(),
		},
		{
			UserID:          uuid.New(),
			CategoryID:      uuid.New(),
			CommunityID:     uuid.New(),
			Content:         "This is my second thread. Excited to be here!",
			ReplyPermission: "everyone",
			IsAdvertisement: false,
			Status:          "published",
			CreatedAt:       time.Now(),
		},
	}

	for _, element := range threads {
		db.Create(&element)
	}
}