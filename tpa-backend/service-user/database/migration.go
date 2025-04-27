package database

import (
	"github.com/Wenfuuu/tpa-web/service-user/model"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.Migrator().DropTable(&model.User{})
	db.Migrator().DropTable(&model.SecurityQuestion{})

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.SecurityQuestion{})
}