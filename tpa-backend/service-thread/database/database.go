package database

import (
	"github.com/Wenfuuu/tpa-web/service-thread/utils"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var database *gorm.DB

func InitDB() *gorm.DB {
	if database == nil {
		host := utils.GetDotEnv("DATABASE_HOST")
		user := utils.GetDotEnv("DATABASE_USER")
		password := utils.GetDotEnv("DATABASE_PASSWORD")
		dbname := utils.GetDotEnv("DATABASE_NAME")
		port := utils.GetDotEnv("DATABASE_PORT")
		sslmode := utils.GetDotEnv("DATABASE_SSLMODE")
		timezone := utils.GetDotEnv("DATABASE_TIMEZONE")

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", host, user, password, dbname, port, sslmode, timezone)
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

		if err != nil {
			panic(err)
		}

		database = db
	}

	return database
}