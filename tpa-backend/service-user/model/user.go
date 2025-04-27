package model

import (
	"time"

	// "github.com/Wenfuuu/tpa-web/service-user/utils"

	// "github.com/Wenfuuu/tpa-web/service-user/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID       			uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Name     			string    `json:"name"`
	Username			string    `json:"username"`
	Email    			string    `json:"email"`
	Password 			string    `json:"password"`
	Gender 				string    `json:"gender"`
	ProfilePicture 		string    `json:"profilePicture"`
	Banner 				string    `json:"banner"`
	DateOfBirth 		string    `json:"dateOfBirth"`
	IsActive        	bool      `json:"isActive" gorm:"default:true"`
	IsBanned        	bool      `json:"isBanned" gorm:"default:false"`
	Role            	string    `json:"role" gorm:"default:'user'"`
	IsSubscribed    	bool      `json:"isSubscribed"`
	IsPremium       	bool      `json:"isPremium" gorm:"default:false"`
	Bio             	string    `json:"bio"`
	CreatedAt        	time.Time `json:"createdAt" gorm:"autoCreateTime"`
}

type SecurityQuestion struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID   uuid.UUID `gorm:"type:uuid;not null" json:"userId"`
	User     User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Question string    `json:"question"`
	Answer   string    `json:"answer"`
}

// BeforeCreate hook will be called automatically by GORM
func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = uuid.New()

	// hashedPassword, err := utils.EncryptPassword(user.Password)
	// if err != nil {
	// 	return err
	// }
	// user.Password = hashedPassword

	return
}
