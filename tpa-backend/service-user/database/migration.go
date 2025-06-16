package database

import (
	"github.com/Wenfuuu/tpa-web/service-user/model"
	"github.com/Wenfuuu/tpa-web/service-user/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	logger := utils.GetLogger()
	logger.Info("Starting database migration")

	// db.Migrator().DropTable(&model.User{})
	// db.Migrator().DropTable(&model.SecurityQuestion{})

	logger.Debug("Dropping existing tables")
	db.Migrator().DropTable(&model.Follow{})
	// db.Migrator().DropTable(&model.PremiumRequest{})
	db.Migrator().DropTable(&model.Report{})

	logger.Debug("Running auto migration for all models")
	if err := db.AutoMigrate(&model.User{}); err != nil {
		logger.Error("Failed to migrate User model", zap.Error(err))
	}
	if err := db.AutoMigrate(&model.SecurityQuestion{}); err != nil {
		logger.Error("Failed to migrate SecurityQuestion model", zap.Error(err))
	}
	if err := db.AutoMigrate(&model.Follow{}); err != nil {
		logger.Error("Failed to migrate Follow model", zap.Error(err))
	}
	if err := db.AutoMigrate(&model.PremiumRequest{}); err != nil {
		logger.Error("Failed to migrate PremiumRequest model", zap.Error(err))
	}
	if err := db.AutoMigrate(&model.Report{}); err != nil {
		logger.Error("Failed to migrate Report model", zap.Error(err))
	}

	logger.Info("Database migration completed successfully")
}
