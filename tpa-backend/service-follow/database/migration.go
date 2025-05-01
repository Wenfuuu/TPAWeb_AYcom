package database

import (
	"github.com/Wenfuuu/tpa-web/service-follow/model"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.Migrator().DropTable(&model.Follow{})

	db.AutoMigrate(&model.Follow{})
}