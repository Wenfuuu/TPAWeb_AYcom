package database

import (
	"github.com/Wenfuuu/tpa-web/service-user/model"
	"github.com/Wenfuuu/tpa-web/service-user/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	logger := utils.GetLogger()
	logger.Info("Starting database seeding")

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

	userId1, err := uuid.Parse("ae2cf219-a6e3-46d7-aee0-6f05eb23f602")
	if err != nil {
		logger.Error("Failed to parse userId1", zap.Error(err))
		return
	}

	userId2, err := uuid.Parse("f3122b31-7b92-400d-b94d-28ef100ca23b")
	if err != nil {
		logger.Error("Failed to parse userId2", zap.Error(err))
		return
	}

	userId3, err := uuid.Parse("812c2aa8-09e0-401d-b0ac-c068f7b93152")
	if err != nil {
		logger.Error("Failed to parse userId3", zap.Error(err))
		return
	}

	userId4, err := uuid.Parse("7e54e720-c237-404c-86d5-e220890db84c")
	if err != nil {
		logger.Error("Failed to parse userId4", zap.Error(err))
		return
	}

	logger.Debug("Seeding follow relationships")
	var follows = []model.Follow{
		{FollowerID: userId1, FollowingID: userId2},
		// {FollowerID: userId1, FollowingID: userId3},
		// {FollowerID: userId1, FollowingID: userId4},
		{FollowerID: userId2, FollowingID: userId1},
		{FollowerID: userId2, FollowingID: userId3},
		{FollowerID: userId2, FollowingID: userId4},
		{FollowerID: userId3, FollowingID: userId1},
		{FollowerID: userId3, FollowingID: userId2},
		// {FollowerID: userId3, FollowingID: userId4},
		{FollowerID: userId4, FollowingID: userId1},
		{FollowerID: userId4, FollowingID: userId2},
		{FollowerID: userId4, FollowingID: userId3},
	}

	for _, element := range follows {
		if err := db.Create(&element).Error; err != nil {
			logger.Warn("Failed to create follow relationship",
				zap.String("follower_id", element.FollowerID.String()),
				zap.String("following_id", element.FollowingID.String()),
				zap.Error(err))
		}
	}

	logger.Debug("Seeding reports")
	var reports = []model.Report{
		{UserID: userId1, ReportedID: userId2, Reason: "Racist"},
		{UserID: userId3, ReportedID: userId4, Reason: "Gay"},
	}

	for _, element := range reports {
		if err := db.Create(&element).Error; err != nil {
			logger.Warn("Failed to create report",
				zap.String("user_id", element.UserID.String()),
				zap.String("reported_id", element.ReportedID.String()),
				zap.Error(err))
		}
	}

	logger.Info("Database seeding completed successfully")
}
