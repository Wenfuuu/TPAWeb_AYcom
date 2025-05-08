// project/backend/service-user/repository/user_repository.go
package repository

import (
	"github.com/Wenfuuu/tpa-web/service-user/database"
	"github.com/Wenfuuu/tpa-web/service-user/model"
	"errors"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		DB: database.InitDB(),
	}
}

func (r *UserRepository) GetSecurityQuestion(userID string) (*model.SecurityQuestion, error) {
	var question model.SecurityQuestion
	result := r.DB.Where("user_id = ?", userID).First(&question)
	if result.Error != nil {
		return nil, result.Error
	}
	return &question, nil
}

func (r *UserRepository) GetFollower(userid string) ([]model.Follow, error) {
	var followers []model.Follow
	result := r.DB.Where("following_id = ?", userid).Find(&followers)
	return followers, result.Error
}

func (r *UserRepository) GetFollowing(userid string) ([]model.Follow, error) {
	var followings []model.Follow
	result := r.DB.Where("follower_id = ?", userid).Find(&followings)
	return followings, result.Error
}

func (r *UserRepository) GetAllUsers() ([]model.User, error) {
	var users []model.User
	result := r.DB.Find(&users)
	return users, result.Error
}

func (r *UserRepository) GetWhoToFollow(userid string) ([]model.User, error) {
	var users []model.User

	limit := 3

	result := r.DB.
		Table("follows").
		Select("users.*").
		Joins("JOIN users ON users.id = follows.following_id").
		Where("users.id != ?", userid).
		Group("users.id").
		Order("COUNT(follows.follower_id) DESC").
		Limit(limit).
		Scan(&users)

	return users, result.Error
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	result := r.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) FindByID(id string) (*model.User, error) {
	var user model.User
	result := r.DB.First(&user, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) RegisterPremium(premium *model.PremiumRequest) error {
	return r.DB.Create(premium).Error
}

func (r *UserRepository) CreateSecurityQuestion(securityQuestion *model.SecurityQuestion) error {
	return r.DB.Create(securityQuestion).Error
}

func (r *UserRepository) FollowUser(follow *model.Follow) (error) {
	return r.DB.Create(follow).Error
}

func (r *UserRepository) UnfollowUser(followerid string, followingid string) (error) {
	result := r.DB.Where("follower_id = ? AND following_id = ?", followerid, followingid).Delete(&model.Follow{})
	return result.Error
}

func (r *UserRepository) Create(user *model.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepository) Update(user *model.User) error {
	return r.DB.Save(user).Error
}

func (r *UserRepository) Delete(id string) (int64, error) {
	result := r.DB.Delete(&model.User{}, "id = ?", id)
	return result.RowsAffected, result.Error
}

func (r *UserRepository) EmailExists(email string) (bool, error) {
	var user model.User
	result := r.DB.Where("email = ?", email).First(&user)
	if result.Error == nil {
		return true, nil
	}
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, result.Error
}

func (r *UserRepository) UsernameExists(username string) (bool, error) {
	var user model.User
	result := r.DB.Where("username = ?", username).First(&user)
	if result.Error == nil {
		return true, nil
	}
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, result.Error
}

func (r *UserRepository) AlreadyRegisteredPremium(user_id string) (bool, error) {
	var request model.PremiumRequest
	result := r.DB.Where("user_id = ?", user_id).First(&request)
	if result.Error == nil {
		return true, nil
	}
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, result.Error
}

func (r *UserRepository) AlreadyPremium(user_id string) (bool, error) {
	var user model.User
	result := r.DB.First(&user, "id = ?", user_id)
	
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return false, nil
        }
        return false, result.Error
    }
    
    return user.IsPremium, nil
}