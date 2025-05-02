package model

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ThreadFiles struct {
	ID       	uuid.UUID 	`gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	ThreadID   	uuid.UUID 	`gorm:"type:uuid;not null" json:"threadId"`
	Thread     	Thread      `gorm:"foreignKey:ThreadID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	File		string		`json:"file"`
}

type ThreadCategory struct {
	ID 		uuid.UUID 	`gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name	string		`json:"name"`
}

type Community struct {
	ID 				uuid.UUID 	`gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name			string		`json:"name"`
	Description    	string    	`gorm:"type:text;not null" json:"description"`
}

type Like struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_thread" json:"userId"`
	ThreadID  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_thread" json:"threadId"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

type Bookmark struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_thread" json:"userId"`
	ThreadID  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_thread" json:"threadId"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

type Thread struct {
	ID             uuid.UUID 	`gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID         uuid.UUID 	`gorm:"type:uuid;not null" json:"userId"`
	CategoryID     uuid.UUID 	`gorm:"type:uuid;not null" json:"categoryId"`
	CommunityID    uuid.UUID 	`gorm:"type:uuid;" json:"communityId"`
	Content        string    	`gorm:"type:text;not null" json:"content"`
	ReplyPermission string   	`gorm:"type:varchar(50);not null" json:"replyPermission"`
	IsAdvertisement bool     	`gorm:"default:false" json:"isAdvertisement"`
	Status         string    	`gorm:"type:varchar(20);default:'published'" json:"status"` // "published" or "scheduled"
	ScheduledAt    *time.Time 	`json:"scheduledAt"` // Pointer to allow null
	CreatedAt      time.Time 	`gorm:"autoCreateTime" json:"createdAt"`
}

func (thread *Thread) BeforeCreate(tx *gorm.DB) (err error) {
	thread.ID = uuid.New()

	if thread.ScheduledAt == nil {
		thread.Status = "published"
	}else{
		thread.Status = "scheduled"
	}

	return
}

func (threadFile *ThreadFiles) BeforeCreate(tx *gorm.DB) (err error) {
	threadFile.ID = uuid.New()
	return
}

func (category *ThreadCategory) BeforeCreate(tx *gorm.DB) (err error) {
	category.ID = uuid.New()
	return
}