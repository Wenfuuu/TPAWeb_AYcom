package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Client represents a connected user
type Client struct {
	Conn          	*websocket.Conn
	UserID        	uuid.UUID 
	ConversationID 	uuid.UUID 
}

type CreateConversationRequest struct {
	Name      string 		`json:"name"`
	UserIDs   []uuid.UUID 	`json:"user_ids"`
	CreatorID uint   		`json:"creator_id"`
}

// Database models
type Conversation struct {
	ID           uuid.UUID 					`gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name         string                   	`json:"name"`
	CreatedAt    time.Time                 	`json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time                 	`json:"updated_at"`
	Participants []ConversationParticipant 	`json:"participants,omitempty" gorm:"foreignKey:ConversationID"`
	Messages     []Message                 	`json:"messages,omitempty" gorm:"foreignKey:ConversationID"`
	IsDeleted 	 bool      					`json:"is_deleted" gorm:"default:false"` // For soft delete
}

type ConversationParticipant struct {
	ID             uuid.UUID 		`gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	ConversationID uuid.UUID      	`json:"conversation_id"`
	UserID         uuid.UUID      	`json:"user_id"`
	JoinedAt       time.Time 		`json:"joined_at" gorm:"autoCreateTime"`
	LastReadAt     time.Time 		`json:"last_read_at"`
	IsActive       bool      		`json:"is_active" gorm:"default:true"` // For removing participants
}

type Message struct {
	ID             uuid.UUID 		`gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	ConversationID uuid.UUID      	`json:"conversation_id"`
	SenderID       uuid.UUID      	`json:"sender_id"`
	Content        string    		`json:"content"`
	MediaType      string    		`json:"media_type"` // "text", "image", "gif", "video"
	MediaURL       string    		`json:"media_url"`  // URL or base64 for media content
	CreatedAt      time.Time 		`json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time 		`json:"updated_at"`
	IsDeleted      bool      		`json:"is_deleted" gorm:"default:false"`
	IsUnsent       bool      		`json:"is_unsent" gorm:"default:false"` // For message unsend feature
}