package database

import (
	"github.com/Wenfuuu/tpa-web/service-user/model"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	// db.Migrator().DropTable(&model.User{})
	// db.Migrator().DropTable(&model.SecurityQuestion{})
	db.Migrator().DropTable(&model.Follow{})

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.SecurityQuestion{})
	db.AutoMigrate(&model.Follow{})
}