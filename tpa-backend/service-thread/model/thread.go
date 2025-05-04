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
	Name	string		`gorm:"unique" json:"name"`
}

type Community struct {
	ID 				uuid.UUID 	`gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name			string		`json:"name"`
	Description    	string    	`gorm:"type:text;not null" json:"description"`
	Logo    		string    	`gorm:"type:text;not null" json:"logo"`
	CategoryID     	uuid.UUID 	`gorm:"type:uuid" json:"categoryId"`
	Banner    		string    	`gorm:"type:text;not null" json:"banner"`
	Rule    		string    	`gorm:"type:text;not null" json:"rule"`
}

type CommunityCategory struct {
	ID 		uuid.UUID 	`gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name	string		`json:"name"`
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
	UserID         uuid.UUID 	`gorm:"type:uuid;not null;uniqueIndex:idx_repost" json:"userId"`
	CategoryID     uuid.UUID 	`gorm:"type:uuid" json:"categoryId"`
	CommunityID    uuid.UUID 	`gorm:"type:uuid;" json:"communityId"`
	ReplyToID      *uuid.UUID 	`gorm:"type:uuid" json:"replyToId"`
	RepostOfID     *uuid.UUID 	`gorm:"type:uuid;uniqueIndex:idx_repost" json:"repostOfId"`
	Content        string    	`gorm:"type:text;not null" json:"content"`
	ReplyPermission string   	`gorm:"type:varchar(50);not null" json:"replyPermission"`
	IsAdvertisement bool     	`gorm:"default:false" json:"isAdvertisement"`
	Status         string    	`gorm:"type:varchar(20);default:'published'" json:"status"` // "published" or "scheduled"
	ScheduledAt    *time.Time 	`json:"scheduledAt"`
	CreatedAt      time.Time 	`gorm:"autoCreateTime" json:"createdAt"`
}

type Hashtag struct {
	ID 		uuid.UUID 	`gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name	string		`json:"name"`
}

type ThreadHashtag struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	ThreadID  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_hashtag" json:"threadId"`
	HashtagID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_hashtag" json:"hashtagId"`
}

type HashtagCount struct {
	ID         	uuid.UUID `json:"id"`
	Name       	string    `json:"name"`
	Count 		int       `json:"count"`
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

func (hashtag *Hashtag) BeforeCreate(tx *gorm.DB) (err error) {
	hashtag.ID = uuid.New()
	return
}

func (threadHashtag *ThreadHashtag) BeforeCreate(tx *gorm.DB) (err error) {
	threadHashtag.ID = uuid.New()
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

func (category *CommunityCategory) BeforeCreate(tx *gorm.DB) (err error) {
	category.ID = uuid.New()
	return
}