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

func (r *ThreadRepository) GetReplies(page int, id string) ([]model.Thread, error) {
	var threads []model.Thread

	limit := 5
	offset := (page - 1) * limit

	result := r.DB.Where("reply_to_id = ?", id).Order("created_at DESC").Limit(limit).Offset(offset).Find(&threads)
	return threads, result.Error
}

func (r *ThreadRepository) GetThreadById(id string) (*model.Thread, error) {
	var thread model.Thread
	result := r.DB.First(&thread, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &thread, nil
}

func (r *ThreadRepository) GetAllThreads(page int) ([]model.Thread, error) {
	var threads []model.Thread

	// add offset here
	limit := 5
	offset := (page - 1) * limit

	result := r.DB.Where("reply_to_id IS NULL").Order("created_at DESC").Limit(limit).Offset(offset).Find(&threads)
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

func (r *ThreadRepository) GetBookmark(page int, userId string) ([]model.Thread, error) {
	var threads []model.Thread

	limit := 5
	offset := (page - 1) * limit

	result := r.DB.
		Joins("JOIN bookmarks ON bookmarks.thread_id = threads.id").
		Where("bookmarks.user_id = ?", userId).
		Offset(offset).
		Find(&threads)
	
	return threads, result.Error
}

func (r *ThreadRepository) GetRepostThreads(userId string) ([]model.Thread, error) {
	var threads []model.Thread

	result := r.DB.
		Where("repost_of_id IS NOT NULL AND user_id = ?", userId).
		Find(&threads)
	
	return threads, result.Error
}

func (r *ThreadRepository) GetLikedThreads(userId string) ([]model.Thread, error) {
	var threads []model.Thread

	result := r.DB.
		Joins("JOIN likes ON likes.thread_id = threads.id").
		Where("likes.user_id = ?", userId).
		Find(&threads)
	
	return threads, result.Error
}

func (r *ThreadRepository) GetBookmarkThreads(userId string) ([]model.Thread, error) {
	var threads []model.Thread

	result := r.DB.
		Joins("JOIN bookmarks ON bookmarks.thread_id = threads.id").
		Where("bookmarks.user_id = ?", userId).
		Find(&threads)
	
	return threads, result.Error
}

func (r *ThreadRepository) UnrepostThread(thread *model.Thread) error {
	var existingThread model.Thread
	if err := r.DB.Where("user_id = ? AND repost_of_id = ?", thread.UserID, thread.RepostOfID).First(&existingThread).Error; err != nil {
		return err
	}

	return r.DB.Delete(&existingThread).Error
}

func (r *ThreadRepository) LikeThread(like *model.Like) error {
	return r.DB.Create(like).Error
}

func (r *ThreadRepository) UnlikeThread(unlike *model.Like) error {
	var existingLike model.Like
	if err := r.DB.Where("user_id = ? AND thread_id = ?", unlike.UserID, unlike.ThreadID).First(&existingLike).Error; err != nil {
		return err
	}

	return r.DB.Delete(&existingLike).Error
}

func (r *ThreadRepository) BookmarkThread(bookmark *model.Bookmark) error {
	return r.DB.Create(bookmark).Error
}

func (r *ThreadRepository) UnbookmarkThread(unbookmark *model.Bookmark) error {
	var existingBookmark model.Bookmark
	if err := r.DB.Where("user_id = ? AND thread_id = ?", unbookmark.UserID, unbookmark.ThreadID).First(&existingBookmark).Error; err != nil {
		return err
	}

	return r.DB.Delete(&existingBookmark).Error
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

func (r *ThreadRepository) GetFileById(thread_id string) ([]model.ThreadFiles, error) {
	var files []model.ThreadFiles
	result := r.DB.Where("thread_id = ?", thread_id).Find(&files)
	return files, result.Error
}

func (r *ThreadRepository) Create(thread *model.Thread) error {
	return r.DB.Create(thread).Error
}

func (r *ThreadRepository) InsertFile(files []model.ThreadFiles) error {
	return r.DB.Create(files).Error
} 