package repository

import (
	"strings"

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

func (r *ThreadRepository) JoinCommunity(req *model.CommunityMembers, userid string, communityid string) error {
	var joinRequest model.JoinCommunityRequest
	if err := r.DB.Where("user_id = ? AND community_id = ?", userid, communityid).First(&joinRequest).Error; err != nil {
		return err
	}

	if err := r.DB.Delete(&joinRequest).Error; err != nil {
		return err
	}

	return r.DB.Create(req).Error
}

func (r *ThreadRepository) GetJoinCommunity(community_id string) ([]model.JoinCommunityRequest, error) {
	var requests []model.JoinCommunityRequest

	result := r.DB.Where("community_id = ?", community_id).Find(&requests)

	return requests, result.Error
}

func (r *ThreadRepository) RequestJoinCommunity(req *model.JoinCommunityRequest) error {
	return r.DB.Create(req).Error
}

func (r *ThreadRepository) GetTopHashtags() ([]model.HashtagCount, error) {
	var hashtags []model.HashtagCount

	limit := 5

	result := r.DB.
		Table("thread_hashtags").
		Select("hashtags.id, hashtags.name, COUNT(thread_hashtags.hashtag_id) as count").
		Joins("JOIN hashtags ON hashtags.id = thread_hashtags.hashtag_id").
		Group("hashtags.id, hashtags.name").
		Order("count DESC").
		Limit(limit).
		Scan(&hashtags)

	return hashtags, result.Error
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

func (r *ThreadRepository) GetTopCommunityThreads(page int, communityId string) ([]model.Thread, error) {
	var threads []model.Thread

	limit := 5
	offset := (page - 1) * limit

	result := r.DB.
		Table("threads").
		Select("threads.*, COUNT(likes.id) AS like_count").
		Joins("LEFT JOIN likes ON likes.thread_id = threads.id").
		Where("threads.community_id = ?", communityId).
		Group("threads.id").
		Order("like_count DESC").
		Limit(limit).
		Offset(offset).
		Scan(&threads)

	return threads, result.Error
}

func (r *ThreadRepository) GetCommunityThreads(page int, communityId string) ([]model.Thread, error) {
	var threads []model.Thread

	// add offset here
	limit := 5
	offset := (page - 1) * limit

	result := r.DB.Where("community_id = ?", communityId).Order("created_at DESC").Limit(limit).Offset(offset).Find(&threads)
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

func (r *ThreadRepository) GetBookmark(page int, userId string, search string) ([]model.Thread, error) {
	var threads []model.Thread

	limit := 5
	offset := (page - 1) * limit

	result := r.DB.
		Joins("JOIN bookmarks ON bookmarks.thread_id = threads.id").
		Where("bookmarks.user_id = ?", userId).
		Where("LOWER(threads.content) LIKE ?", "%"+strings.ToLower(search)+"%").
		Limit(limit).
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

func (r *ThreadRepository) GetCommunityById(id string) (*model.Community, error) {
	var community model.Community
	result := r.DB.First(&community, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &community, nil
}

func (r *ThreadRepository) GetCommunityMembers(community_id string) ([]model.CommunityMembers, error) {
	var members []model.CommunityMembers

	result := r.DB.
		Where("community_id = ?", community_id).
		Find(&members)

	return members, result.Error
}

func (r *ThreadRepository) GetTopCommunityMembers(community_id string) ([]model.CommunityMembers, error) {
	var members []model.CommunityMembers

	result := r.DB.
		Table("community_members").
		Select("community_members.id, community_members.user_id, community_members.community_id, community_members.role, COUNT(follows.follower_id) AS follower_count").
		Joins("LEFT JOIN follows ON follows.following_id = community_members.user_id").
		Where("community_members.community_id = ?", community_id).
		Group("community_members.id, community_members.user_id, community_members.community_id, community_members.role").
		Order("follower_count DESC").
		Limit(3).
		Find(&members)

	return members, result.Error
}

func (r *ThreadRepository) GetAllCommunity() ([]model.Community, error) {
	var communities []model.Community
	result := r.DB.Find(&communities)
	return communities, result.Error
}

func (r *ThreadRepository) GetPendingCommunity(userid string, page int, limit int, search string, category_id string) ([]model.Community, error) {
	var communities []model.Community

	offset := (page - 1) * limit

	query := r.DB.
		Joins("JOIN join_community_requests ON join_community_requests.community_id = communities.id").
		Where("join_community_requests.user_id = ?", userid).
		Where("LOWER(communities.name) LIKE ?", "%"+strings.ToLower(search)+"%")

	if category_id != "all" {
		query = query.Where("communities.category_id = ?", category_id)
	}

	result := query.
		Limit(limit).
		Offset(offset).
		Find(&communities)

	return communities, result.Error
}

func (r *ThreadRepository) GetJoinedCommunity(userid string, page int, limit int, search string, category_id string) ([]model.Community, error) {
	var communities []model.Community

	offset := (page - 1) * limit

	query := r.DB.
		Joins("JOIN community_members ON community_members.community_id = communities.id").
		Where("community_members.user_id = ?", userid).
		Where("LOWER(communities.name) LIKE ?", "%"+strings.ToLower(search)+"%")

	if category_id != "all" {
		query = query.Where("communities.category_id = ?", category_id)
	}

	result := query.
		Limit(limit).
		Offset(offset).
		Find(&communities)

	return communities, result.Error
}

func (r *ThreadRepository) GetCommunityPagination(page int, limit int, search string, category_id string) ([]model.Community, error) {
	var communities []model.Community

	offset := (page - 1) * limit

	query := r.DB.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(search)+"%")

	if category_id != "all" {
		query = query.Where("category_id = ?", category_id)
	}

	result := query.
		Limit(limit).
		Offset(offset).
		Find(&communities)

	return communities, result.Error
}

func (r *ThreadRepository) GetAllCommunityCategories() ([]model.CommunityCategory, error) {
	var categories []model.CommunityCategory
	result := r.DB.Find(&categories)
	return categories, result.Error
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
