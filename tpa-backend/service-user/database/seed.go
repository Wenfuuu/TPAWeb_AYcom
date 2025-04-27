package database

import (
	"github.com/Wenfuuu/tpa-web/service-user/model"

	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	var users = []model.User{
		{
			Name:           "Alice Smith",
			Username:       "alice123",
			Email:          "email1@gmail.com",
			Password:       "pass1",
			Gender:         "Female",
			ProfilePicture: "https://example.com/profiles/alice.jpg",
			Banner:         "https://example.com/banners/alice.jpg",
			DateOfBirth:    "1995-06-15",
			IsActive:       true,
			IsBanned:       false,
			Role:           "user",
			IsSubscribed: 	false,
			IsPremium: 		false,
			Bio: 			"hi niggers",
		},
		{
			Name:           "Carol Lee",
			Username:       "carollee",
			Email:          "email2@gmail.com",
			Password:       "pass2",
			Gender:         "Female",
			ProfilePicture: "https://example.com/profiles/carol.jpg",
			Banner:         "https://example.com/banners/carol.jpg",
			DateOfBirth:    "1992-03-22",
			IsActive:       true,
			IsBanned:       true,
			Role:           "user",
			IsSubscribed: 	true,
			IsPremium: 		true,
			Bio: 			"maklo",
		},
	}

	for _, element := range users {
		db.Create(&element)
	}
}