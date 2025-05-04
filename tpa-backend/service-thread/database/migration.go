package database

import (
	"github.com/Wenfuuu/tpa-web/service-thread/model"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	// db.Migrator().DropTable(&model.Thread{})
	// db.Migrator().DropTable(&model.ThreadFiles{})
	// db.Migrator().DropTable(&model.ThreadCategory{})
	// db.Migrator().DropTable(&model.CommunityCategory{})
	// db.Migrator().DropTable(&model.Community{})
	// db.Migrator().DropTable(&model.Like{})
	// db.Migrator().DropTable(&model.Bookmark{})
	// db.Migrator().DropTable(&model.Hashtag{})
	// db.Migrator().DropTable(&model.ThreadHashtag{})

	db.AutoMigrate(&model.Thread{})
	db.AutoMigrate(&model.ThreadFiles{})
	db.AutoMigrate(&model.ThreadCategory{})
	db.AutoMigrate(&model.Community{})
	db.AutoMigrate(&model.CommunityCategory{})
	db.AutoMigrate(&model.Like{})
	db.AutoMigrate(&model.Bookmark{})
	db.AutoMigrate(&model.Hashtag{})
	db.AutoMigrate(&model.ThreadHashtag{})
}