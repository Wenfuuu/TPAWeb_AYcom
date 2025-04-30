package repository

import (
	"github.com/Wenfuuu/tpa-web/service-thread/database"
	"github.com/Wenfuuu/tpa-web/service-thread/model"
	"gorm.io/gorm"
)

type ThreadRepository struct {
	DB *gorm.DB
}

func NewThreadRepository() *ThreadRepository {
	return &ThreadRepository{
		DB: database.InitDB(),
	}
}

func (r *ThreadRepository) GetAllThreads() ([]model.Thread, error) {
	var threads []model.Thread
	result := r.DB.Find(&threads)
	return threads, result.Error
}