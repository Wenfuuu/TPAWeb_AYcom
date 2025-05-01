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

func (r *ThreadRepository) GetAllThreads(page int) ([]model.Thread, error) {
	var threads []model.Thread

	// add offset here
	limit := 5
	offset := (page - 1) * limit

	result := r.DB.Order("created_at DESC").Limit(limit).Offset(offset).Find(&threads)
	return threads, result.Error
}

func (r *ThreadRepository) GetFollowingThreads(page int, userIds []string) ([]model.Thread, error) {
	var threads []model.Thread

	// add offset here
	limit := 5
	offset := (page - 1) * limit

	result := r.DB.Where("user_id IN ?", userIds).Order("created_at DESC").Limit(limit).Offset(offset).Find(&threads)
	return threads, result.Error
}

func (r *ThreadRepository) GetAllThreadCategories() ([]model.ThreadCategory, error) {
	var categories []model.ThreadCategory
	result := r.DB.Find(&categories)
	return categories, result.Error
}

func (r *ThreadRepository) GetAllCommunity() ([]model.Community, error) {
	var communities []model.Community
	result := r.DB.Find(&communities)
	return communities, result.Error
}

func (r *ThreadRepository) GetAllThreadFile() ([]model.ThreadFiles, error) {
	var files []model.ThreadFiles
	result := r.DB.Find(&files)
	return files, result.Error
}

func (r *ThreadRepository) Create(thread *model.Thread) error {
	return r.DB.Create(thread).Error
}

func (r *ThreadRepository) InsertFile(files []model.ThreadFiles) error {
	return r.DB.Create(files).Error
} 