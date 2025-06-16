package database

import (
	"fmt"

	"github.com/Wenfuuu/tpa-web/service-user/utils"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var database *gorm.DB

func InitDB() *gorm.DB {
	logger := utils.GetLogger()

	if database == nil {
		logger.Info("Initializing database connection")

		host := utils.GetDotEnv("DATABASE_HOST")
		user := utils.GetDotEnv("DATABASE_USER")
		password := utils.GetDotEnv("DATABASE_PASSWORD")
		dbname := utils.GetDotEnv("DATABASE_NAME")
		port := utils.GetDotEnv("DATABASE_PORT")
		sslmode := utils.GetDotEnv("DATABASE_SSLMODE")
		timezone := utils.GetDotEnv("DATABASE_TIMEZONE")

		logger.Debug("Database configuration loaded",
			zap.String("host", host),
			zap.String("user", user),
			zap.String("dbname", dbname),
			zap.String("port", port),
			zap.String("sslmode", sslmode),
			zap.String("timezone", timezone))

		// dsn := "host=localhost user=postgres password=root dbname=tpa-web port=5432 sslmode=disable TimeZone=Asia/Jakarta"
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", host, user, password, dbname, port, sslmode, timezone)
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

		if err != nil {
			logger.Fatal("Failed to connect to database", zap.Error(err))
		}

		logger.Info("Database connection established successfully")
		database = db
	}

	return database
}
