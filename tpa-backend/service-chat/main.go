package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/Wenfuuu/tpa-web/pkg/messaging"
	chatpb "github.com/Wenfuuu/tpa-web/pkg/proto/chat"
	"github.com/Wenfuuu/tpa-web/service-chat/database"
	"github.com/Wenfuuu/tpa-web/service-chat/model"
	"github.com/Wenfuuu/tpa-web/service-chat/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketMessage struct {
	Type    string      `json:"type"`    // "chat" or "notification"
	Payload interface{} `json:"payload"` // The actual message data
}

var db *gorm.DB

// Global map of active clients
var clients = make(map[*model.Client]bool)

// Channel to broadcast messages
var broadcast = make(chan model.Message)

// Channel for notifications
var notifyChannel = make(chan model.Notification)

var rabbitChannel *amqp091.Channel

type ChatServer struct {
	chatpb.UnimplementedChatServiceServer
}

func handleMessages() {
	for {
		select {
		case msg := <-broadcast:
			// Handle chat messages as before
			for client := range clients {
				wsMessage := WebSocketMessage{
					Type:    "chat",
					Payload: msg,
				}
				err := client.Conn.WriteJSON(wsMessage)
				if err != nil {
					log.Printf("Error sending message: %v", err)
					client.Conn.Close()
					delete(clients, client)
				}
			}

		case notification := <-notifyChannel:
			// Send notification to the specific user
			for client := range clients {
				if client.UserID == notification.UserID {
					wsMessage := WebSocketMessage{
						Type:    "notification",
						Payload: notification,
					}
					err := client.Conn.WriteJSON(wsMessage)
					if err != nil {
						log.Printf("Error sending notification: %v", err)
						client.Conn.Close()
						delete(clients, client)
					}
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

	// Parse user ID - this is always required
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Printf("Invalid user ID: %v", err)
		return
	}

	// Create a new client with required user ID
	client := &model.Client{
		Conn:   ws,
		UserID: userID,
	}

	// Parse conversation ID only if provided
	if convIDStr != "" {
		convID, err := uuid.Parse(convIDStr)
		if err != nil {
			log.Printf("Invalid conversation ID: %v", err)
			return
		}
		client.ConversationID = convID

		// Update participant's last read timestamp only if we have a conversation
		db.Model(&model.ConversationParticipant{}).
			Where("conversation_id = ? AND user_id = ?", convID, userID).
			Updates(map[string]interface{}{"last_read_at": time.Now()})
	}

	// Register client
	clients[client] = true

	// Handle incoming messages
	for {
		var wsMsg WebSocketMessage

		err := ws.ReadJSON(&wsMsg)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			delete(clients, client)
			break
		}

		switch wsMsg.Type {
		case "chat":
			var msg struct {
				Content   string `json:"content"`
				MediaType string `json:"media_type"`
				MediaURL  string `json:"media_url"`
			}

			payloadBytes, _ := json.Marshal(wsMsg.Payload)
			json.Unmarshal(payloadBytes, &msg)

			// Create a new message
			newMsg := model.Message{
				SenderID:       client.UserID,
				ConversationID: client.ConversationID,
				Content:        msg.Content,
				MediaType:      msg.MediaType,
				MediaURL:       msg.MediaURL,
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

		case "notification":
			var notifData struct {
				RecipientID string `json:"recipient_id"`
				SenderID    string `json:"sender_id"`
				ThreadID    string `json:"thread_id,omitempty"`
				Content     string `json:"content"`
				Category    string `json:"category"`
				Email 		string `json:"email"`
			}

			payloadBytes, _ := json.Marshal(wsMsg.Payload)
			json.Unmarshal(payloadBytes, &notifData)

			recipientID, err := uuid.Parse(notifData.RecipientID)
			if err != nil {
				log.Printf("Invalid recipient ID: %v", err)
				continue
			}

			senderID, err := uuid.Parse(notifData.SenderID)
			if err != nil {
				log.Printf("Invalid sender ID: %v", err)
				continue
			}

			var threadID *uuid.UUID
			if notifData.ThreadID != "" {
				parsedThreadID, err := uuid.Parse(notifData.ThreadID)
				if err != nil {
					log.Printf("Invalid thread ID: %v", err)
					continue
				}
				threadID = &parsedThreadID
			}

			// Create notification
			notification := model.Notification{
				UserID: recipientID,
				SenderID: senderID,
				ThreadID: threadID,
				Content:  notifData.Content,
				Category: notifData.Category,
			}

			// Save to database
			if err := db.Create(&notification).Error; err != nil {
				log.Printf("Error saving notification: %v", err)
				continue
			}

			// Send email
			payload := messaging.EmailPayload{
				Email: notifData.Email,
				Word: notifData.Content,
			}
			body, err := json.Marshal(payload)
			if err != nil {
				log.Printf("Failed to marshal email payload: %v", err)
				continue
			}
			exchangeName := "email_exchange"
			routingKey := "email"
			err = rabbitChannel.PublishWithContext(
				context.Background(),
				exchangeName,
				routingKey,
				false,
				false,
				amqp091.Publishing{
					ContentType: "application/json",
					Body:        body,
					DeliveryMode: amqp091.Transient,
				},
			)
			if err != nil {
				log.Printf("Failed to publish email message: %v", err)
				continue
			}
			log.Printf("Successfully published email request for %s", notifData.Email)

			// Send to notification channel
			notifyChannel <- notification
		}
	}
}

func main() {

	ch, err := utils.GetRabbitMQChannel()
	if err != nil {
		log.Fatalf("Failed to get RabbitMQ channel for AuthHandler: %v", err)
	}
	exchangeName := "email_exchange"
	err = ch.ExchangeDeclare(
		exchangeName, // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	rabbitChannel = ch

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
		chatpb.RegisterChatServiceServer(grpcServer, &ChatServer{})

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
	err = router.Run(":" + httpPort)
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

	result := query.Order("created_at asc").Find(&messages)
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
	fmt.Printf(req.ConversationId)
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

func (s *ChatServer) UpdateParticipant(ctx context.Context, req *chatpb.UpdateParticipantRequest) (*chatpb.Conversation, error) {
	conversationID, err := uuid.Parse(req.ConversationId)
	if err != nil {
		return nil, fmt.Errorf("invalid conversation ID: %w", err)
	}

	tx := db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	// Step 1: Delete all existing participants
	if err := tx.Where("conversation_id = ?", conversationID).
		Delete(&model.ConversationParticipant{}).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to remove old participants: %w", err)
	}

	// Step 2: Add new participants
	now := time.Now()
	for _, userIDStr := range req.UserIds {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("invalid user ID %s: %w", userIDStr, err)
		}

		participant := model.ConversationParticipant{
			ConversationID: conversationID,
			UserID:         userID,
			JoinedAt:       now,
			LastReadAt:     now,
			IsActive:       true,
		}

		if err := tx.Create(&participant).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to add participant %s: %w", userIDStr, err)
		}
	}

	// Step 3: Update conversation name if provided
	if req.NewName != "" {
		if err := tx.Model(&model.Conversation{}).
			Where("id = ?", conversationID).
			Update("name", req.NewName).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to update conversation name: %w", err)
		}
	}

	// Step 4: Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("transaction commit failed: %w", err)
	}

	// Step 5: Fetch updated conversation with participants
	var conversation model.Conversation
	if err := db.Preload("Participants").
		First(&conversation, "id = ?", conversationID).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch updated conversation: %w", err)
	}

	return convertConversationToProto(conversation), nil
}

func (s *ChatServer) GetAllNotifications(ctx context.Context, req *chatpb.GetAllNotificationsRequest) (*chatpb.GetNotificationsResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	var notifications []model.Notification
	query := db.Where("user_id = ?", userID).Where("category = ?", "message");

	// Fetch notifications
	if err := query.Order("created_at desc").Find(&notifications).Error; err != nil {
		return nil, err
	}

	// Convert to protobuf format
	pbNotifications := make([]*chatpb.Notification, len(notifications))
	for i, notif := range notifications {
		threadID := ""
		if notif.ThreadID != nil {
			threadID = notif.ThreadID.String()
		}

		pbNotifications[i] = &chatpb.Notification{
			Id:        notif.ID.String(),
			UserId:    notif.UserID.String(),
			SenderId:  notif.SenderID.String(),
			ThreadId:  threadID,
			Content:   notif.Content,
			Category:  notif.Category,
			CreatedAt: notif.CreatedAt.Format(time.RFC3339),
		}
	}

	return &chatpb.GetNotificationsResponse{
		Notifications: pbNotifications,
	}, nil
}

func (s *ChatServer) GetMentionNotifications(ctx context.Context, req *chatpb.GetMentionNotificationsRequest) (*chatpb.GetNotificationsResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	var notifications []model.Notification
	query := db.Where("user_id = ?", userID).Where("category = ?", "mention");

	// Fetch notifications
	if err := query.Order("created_at desc").Find(&notifications).Error; err != nil {
		return nil, err
	}

	// Convert to protobuf format
	pbNotifications := make([]*chatpb.Notification, len(notifications))
	for i, notif := range notifications {
		threadID := ""
		if notif.ThreadID != nil {
			threadID = notif.ThreadID.String()
		}
		
		pbNotifications[i] = &chatpb.Notification{
			Id:        notif.ID.String(),
			UserId:    notif.UserID.String(),
			SenderId:  notif.SenderID.String(),
			ThreadId:  threadID,
			Content:   notif.Content,
			Category:  notif.Category,
			CreatedAt: notif.CreatedAt.Format(time.RFC3339),
		}
	}

	return &chatpb.GetNotificationsResponse{
		Notifications: pbNotifications,
	}, nil
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
		Id:             msg.ID.String(),
		ConversationId: msg.ConversationID.String(),
		SenderId:       msg.SenderID.String(),
		Content:        msg.Content,
		MediaType:      msg.MediaType,
		MediaUrl:       msg.MediaURL,
		CreatedAt:      msg.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      msg.UpdatedAt.Format(time.RFC3339),
		IsDeleted:      msg.IsDeleted,
	}
}
