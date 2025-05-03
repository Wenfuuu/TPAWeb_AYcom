package database

import (
	"github.com/Wenfuuu/tpa-web/service-user/model"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	// var users = []model.User{
	// 	{
	// 		Name:           "Alice Smith",
	// 		Username:       "alice123",
	// 		Email:          "email1@gmail.com",
	// 		Password:       "pass1",
	// 		Gender:         "Female",
	// 		ProfilePicture: "https://example.com/profiles/alice.jpg",
	// 		Banner:         "https://example.com/banners/alice.jpg",
	// 		DateOfBirth:    "1995-06-15",
	// 		IsActive:       true,
	// 		IsBanned:       false,
	// 		Role:           "user",
	// 		IsSubscribed: 	false,
	// 		IsPremium: 		false,
	// 		Bio: 			"hi niggers",
	// 	},
	// 	{
	// 		Name:           "Carol Lee",
	// 		Username:       "carollee",
	// 		Email:          "email2@gmail.com",
	// 		Password:       "pass2",
	// 		Gender:         "Female",
	// 		ProfilePicture: "https://example.com/profiles/carol.jpg",
	// 		Banner:         "https://example.com/banners/carol.jpg",
	// 		DateOfBirth:    "1992-03-22",
	// 		IsActive:       true,
	// 		IsBanned:       true,
	// 		Role:           "user",
	// 		IsSubscribed: 	true,
	// 		IsPremium: 		true,
	// 		Bio: 			"maklo",
	// 	},
	// }

	// for _, element := range users {
	// 	db.Create(&element)
	// }

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