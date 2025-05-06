package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	chatpb "github.com/Wenfuuu/tpa-web/pkg/proto/chat"
	"github.com/Wenfuuu/tpa-web/service-chat/database"
	"github.com/Wenfuuu/tpa-web/service-chat/model"
	"github.com/Wenfuuu/tpa-web/service-chat/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections for now, should restrict in production
	},
}

var db *gorm.DB
// Global map of active clients
var clients = make(map[*model.Client]bool)

// Channel to broadcast messages
var broadcast = make(chan model.Message)

type ChatServer struct {
	chatpb.UnimplementedChatServiceServer
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		// Send it to all clients in the same conversation
		for client := range clients {
			if client.ConversationID == msg.ConversationID {
				err := client.Conn.WriteJSON(msg)
				if err != nil {
					log.Printf("Error sending message: %v", err)
					client.Conn.Close()
					delete(clients, client)
				}
			}
		}
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to websocket: %v", err)
		return
	}
	defer ws.Close()

	// Extract user ID and conversation ID from query params
	userIDStr := r.URL.Query().Get("user_id")
	convIDStr := r.URL.Query().Get("conversation_id")
	
	// Parse user ID and conversation ID
	userID := uuid.MustParse(userIDStr)
	convID := uuid.MustParse(convIDStr)
	
	// In production, validate and parse these values properly
	
	// Create a new client
	client := &model.Client{
		Conn:           ws,
		UserID:         userID,
		ConversationID: convID,
	}
	
	// Register client
	clients[client] = true

	// Update participant's last read timestamp
	db.Model(&model.ConversationParticipant{}).
		Where("conversation_id = ? AND user_id = ?", convID, userID).
		Updates(map[string]interface{}{"last_read_at": time.Now()})

	// Handle incoming messages
	for {
		var msg struct {
			Content   string `json:"content"`
			MediaType string `json:"media_type"`
			MediaURL string `json:"media_url"` // url for images/videos
		}
		
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			delete(clients, client)
			break
		}
		
		// Create a new message
		newMsg := model.Message{
			SenderID:      client.UserID,
			ConversationID: client.ConversationID,
			Content:       msg.Content,
			MediaType:     msg.MediaType,
		}
		
		// Handle media if present
		if msg.MediaType != "text" && msg.MediaURL != "" {
			// In production, you would:
			// 1. Decode base64 data
			// 2. Validate media type
			// 3. Store in blob storage (S3, etc.)
			// 4. Store the URL in the database
			
			// For now, just store a placeholder URL or the base64 data
			// if strings.HasPrefix(msg.MediaURL, "data:") {
			// 	newMsg.MediaURL = msg.MediaURL // Already formatted as data URL
			// } else {
			// 	// Create a data URL
			// 	mediaPrefix := "data:" + msg.MediaType + ";base64,"
			// 	newMsg.MediaURL = mediaPrefix + msg.MediaURL
			// }
		}
		
		// Save message to database
		if err := db.Create(&newMsg).Error; err != nil {
			log.Printf("Error saving message: %v", err)
			continue
		}
		
		// Update conversation's updated_at timestamp
		db.Model(&model.Conversation{}).Where("id = ?", client.ConversationID).
			Updates(map[string]interface{}{"updated_at": time.Now()})
		
		// Update sender's last read timestamp
		db.Model(&model.ConversationParticipant{}).
			Where("conversation_id = ? AND user_id = ?", client.ConversationID, client.UserID).
			Updates(map[string]interface{}{"last_read_at": time.Now()})
		
		// Broadcast message to all clients in the same conversation
		broadcast <- newMsg
	}
}

func main() {

	db = database.InitDB()
	database.Migrate(db)

	go handleMessages()

	// Start gRPC server in background
	go func() {
		grpcPort := utils.GetDotEnv("GRPC_PORT")
		if grpcPort == "" {
			grpcPort = "50055"
		}

		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}

		grpcServer := grpc.NewServer()
		chatpb.RegisterChatServiceServer(grpcServer, &ChatServer{}) // Implement ChatServer

		log.Printf("Chat service listening on port %s", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Create Gin router for WebSocket and HTTP endpoints
	router := gin.Default()

	// WebSocket endpoint
	router.GET("/ws", func(c *gin.Context) {
		handleConnections(c.Writer, c.Request)
	})

	// Start HTTP server
	httpPort := utils.GetDotEnv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "50056"
	}
	log.Printf("Chat service HTTP server starting on port %s", httpPort)
	err := router.Run(":" + httpPort)
	if err != nil {
		log.Fatal("Error starting HTTP server: ", err)
	}

}

// gRPC implementations
func (s *ChatServer) GetConversations(ctx context.Context, req *chatpb.GetConversationsRequest) (*chatpb.GetConversationsResponse, error) {
	var conversations []model.Conversation
	
	query := db.Preload("Participants").
		Joins("JOIN conversation_participants ON conversations.id = conversation_participants.conversation_id").
		Where("conversation_participants.user_id = ?", req.UserId)
	
	if req.Limit > 0 {
		query = query.Limit(int(req.Limit))
	}
	if req.Offset > 0 {
		query = query.Offset(int(req.Offset))
	}
	
	result := query.Order("conversations.updated_at desc").Find(&conversations)
	if result.Error != nil {
		return nil, result.Error
	}
	
	// Count total
	var totalCount int64
	db.Model(&model.Conversation{}).
		Joins("JOIN conversation_participants ON conversations.id = conversation_participants.conversation_id").
		Where("conversation_participants.user_id = ?", req.UserId).
		Count(&totalCount)
		
	// Convert to protobuf format
	pbConversations := make([]*chatpb.Conversation, len(conversations))
	for i, conv := range conversations {
		pbConversations[i] = convertConversationToProto(conv)
	}
	
	return &chatpb.GetConversationsResponse{
		Conversations: pbConversations,
		TotalCount:    int32(totalCount),
	}, nil
}

func (s *ChatServer) GetConversationById(ctx context.Context, req *chatpb.GetConversationByIdRequest) (*chatpb.Conversation, error) {
	var conversation model.Conversation
	
	if err := db.Preload("Participants").First(&conversation, req.ConversationId).Error; err != nil {
		return nil, err
	}
	
	// Check if the user is a participant
	var participant model.ConversationParticipant
	if err := db.Where("conversation_id = ? AND user_id = ?", req.ConversationId, req.UserId).
		First(&participant).Error; err != nil {
		return nil, err
	}
	
	return convertConversationToProto(conversation), nil
}

func (s *ChatServer) CreateConversation(ctx context.Context, req *chatpb.CreateConversationRequest) (*chatpb.Conversation, error) {
	conversation := model.Conversation{
		Name: req.Name,
	}
	
	// Start a transaction
	tx := db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	
	if err := tx.Create(&conversation).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	
	// Add participants
	for _, userID := range req.UserIds {
		participant := model.ConversationParticipant{
			ConversationID: conversation.ID,
			UserID:         uuid.MustParse(userID),
			JoinedAt:       time.Now(),
			LastReadAt:     time.Now(),
		}
		
		if err := tx.Create(&participant).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	
	// Add creator if not already in participants
	creatorInList := false
	for _, userID := range req.UserIds {
		if userID == req.CreatorId {
			creatorInList = true
			break
		}
	}
	
	if !creatorInList {
		participant := model.ConversationParticipant{
			ConversationID: conversation.ID,
			UserID:         uuid.MustParse(req.CreatorId),
			JoinedAt:       time.Now(),
			LastReadAt:     time.Now(),
		}
		
		if err := tx.Create(&participant).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	
	// Fetch the full conversation with participants
	if err := db.Preload("Participants").First(&conversation, conversation.ID).Error; err != nil {
		return nil, err
	}
	
	return convertConversationToProto(conversation), nil
}

func (s *ChatServer) GetMessages(ctx context.Context, req *chatpb.GetMessagesRequest) (*chatpb.GetMessagesResponse, error) {
	var messages []model.Message
	
	query := db.Where("conversation_id = ? AND is_deleted = ?", req.ConversationId, false)
	
	if req.Limit > 0 {
		query = query.Limit(int(req.Limit))
	}
	if req.Offset > 0 {
		query = query.Offset(int(req.Offset))
	}
	
	result := query.Order("created_at desc").Find(&messages)
	if result.Error != nil {
		return nil, result.Error
	}
	
	// Count total
	var totalCount int64
	db.Model(&model.Message{}).
		Where("conversation_id = ? AND is_deleted = ?", req.ConversationId, false).
		Count(&totalCount)
		
	// Convert to protobuf format
	pbMessages := make([]*chatpb.Message, len(messages))
	for i, msg := range messages {
		pbMessages[i] = convertMessageToProto(msg)
	}
	
	return &chatpb.GetMessagesResponse{
		Messages:   pbMessages,
		TotalCount: int32(totalCount),
	}, nil
}

func (s *ChatServer) SendMessage(ctx context.Context, req *chatpb.SendMessageRequest) (*chatpb.Message, error) {
	msg := model.Message{
		ConversationID: uuid.MustParse(req.ConversationId),
		SenderID:       uuid.MustParse(req.SenderId),
		Content:        req.Content,
	}
	
	if err := db.Create(&msg).Error; err != nil {
		return nil, err
	}
	
	// Update conversation's updated_at timestamp
	db.Model(&model.Conversation{}).Where("id = ?", req.ConversationId).
		Updates(map[string]interface{}{"updated_at": time.Now()})
	
	// Broadcast the message to connected clients
	broadcast <- msg
	
	return convertMessageToProto(msg), nil
}

// Helper functions for converting between database models and protobuf messages
func convertConversationToProto(conv model.Conversation) *chatpb.Conversation {
	participants := make([]*chatpb.ConversationParticipant, len(conv.Participants))
	for i, p := range conv.Participants {
		participants[i] = &chatpb.ConversationParticipant{
			Id:             p.ID.String(),
			ConversationId: p.ConversationID.String(),
			UserId:         p.UserID.String(),
			JoinedAt:       p.JoinedAt.Format(time.RFC3339),
			LastReadAt:     p.LastReadAt.Format(time.RFC3339),
		}
	}

	return &chatpb.Conversation{
		Id:           conv.ID.String(),
		Name:         conv.Name,
		CreatedAt:    conv.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    conv.UpdatedAt.Format(time.RFC3339),
		Participants: participants,
	}
}

func convertMessageToProto(msg model.Message) *chatpb.Message {
	return &chatpb.Message{
		Id:            	msg.ID.String(),
		ConversationId: msg.ConversationID.String(),
		SenderId:      	msg.SenderID.String(),
		Content:       	msg.Content,
		CreatedAt:     	msg.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     	msg.UpdatedAt.Format(time.RFC3339),
		IsDeleted:     	msg.IsDeleted,
	}
}