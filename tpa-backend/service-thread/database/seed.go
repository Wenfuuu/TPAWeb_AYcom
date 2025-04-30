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

	var categories = []model.ThreadCategory{
		{
			Name: "Anime",
		},
		{
			Name: "Manga",
		},
		{
			Name: "Manhwa",
		},
	}
	for _, element := range categories {
		db.Create(&element)
	}

	var communities = []model.Community{
		{
			Name: "Hololive",
			Description: "A very special place for everyone",
		},
		{
			Name: "Hoyoverse",
			Description: "A place that talks about the best game creator",
		},
	}
	for _, element := range communities {
		db.Create(&element)
	}
}