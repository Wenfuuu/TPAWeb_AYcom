package repository

import (
	"github.com/Wenfuuu/tpa-web/service-follow/database"
	"github.com/Wenfuuu/tpa-web/service-follow/model"
	"gorm.io/gorm"
)

type FollowRepository struct {
	DB *gorm.DB
}

func NewFollowRepository() *FollowRepository {
	return &FollowRepository{
		DB: database.InitDB(),
	}
}

func (r *FollowRepository) GetFollower(userid string) ([]model.Follow, error) {
	var followers []model.Follow
	result := r.DB.Where("following_id = ?", userid).Find(&followers)
	return followers, result.Error
}

func (r *FollowRepository) GetFollowing(userid string) ([]model.Follow, error) {
	var followings []model.Follow
	result := r.DB.Where("follower_id = ?", userid).Find(&followings)
	return followings, result.Error
}