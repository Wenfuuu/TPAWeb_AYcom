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
		},
		{
			Name:           "Bob Johnson",
			Username:       "bobbyj",
			Email:          "email2@gmail.com",
			Password:       "pass2",
			Gender:         "Male",
			ProfilePicture: "https://example.com/profiles/bob.jpg",
			Banner:         "https://example.com/banners/bob.jpg",
			DateOfBirth:    "1990-11-30",
			IsActive:       true,
			IsBanned:       false,
			Role:           "user",
		},
		{
			Name:           "Carol Lee",
			Username:       "carollee",
			Email:          "email3@gmail.com",
			Password:       "pass3",
			Gender:         "Female",
			ProfilePicture: "https://example.com/profiles/carol.jpg",
			Banner:         "https://example.com/banners/carol.jpg",
			DateOfBirth:    "1992-03-22",
			IsActive:       true,
			IsBanned:       true,
			Role:           "user",
		},
	}

	for _, element := range users {
		db.Create(&element)
	}
}