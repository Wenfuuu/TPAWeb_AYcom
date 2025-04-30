package database

import (
	"github.com/Wenfuuu/tpa-web/service-thread/model"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.Migrator().DropTable(&model.Thread{})

	db.AutoMigrate(&model.Thread{})
}