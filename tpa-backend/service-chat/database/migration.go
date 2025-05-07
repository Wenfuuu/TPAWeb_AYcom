package database

import (
	"github.com/Wenfuuu/tpa-web/service-chat/model"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	// db.Migrator().DropTable(&model.Notification{})

	db.AutoMigrate(&model.Conversation{})
	db.AutoMigrate(&model.ConversationParticipant{})
	db.AutoMigrate(&model.Message{})
	db.AutoMigrate(&model.Notification{})
}