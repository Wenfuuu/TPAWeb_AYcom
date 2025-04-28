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

func (r *UserRepository) GetAllUsers() ([]model.User, error) {
	var users []model.User
	result := r.DB.Find(&users)
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

func (r *UserRepository) CreateSecurityQuestion(securityQuestion *model.SecurityQuestion) error {
	return r.DB.Create(securityQuestion).Error
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