package handler

import (
	"context"
	"net/http"

	chatpb "github.com/Wenfuuu/tpa-web/pkg/proto/chat"
	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	Client chatpb.ChatServiceClient
}

func NewChatHandler(chatClient chatpb.ChatServiceClient) *ChatHandler {
	return &ChatHandler{
		Client: chatClient,
	}
}

func (h *ChatHandler) GetConversations(c *gin.Context) {
	userID := c.Query("user_id")

	resp, err := h.Client.GetConversations(context.Background(), &chatpb.GetConversationsRequest{
		UserId: userID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"conversations": resp.Conversations,
		"total_count": resp.TotalCount,
	})
}

func (h *ChatHandler) CreateConversation(c *gin.Context) {
	var requestBody struct {
		Name      string   `json:"name" binding:"required"`
		UserIDs   []string `json:"user_ids" binding:"required"` // UUIDs in string format
		CreatorID string   `json:"creator_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.Client.CreateConversation(context.Background(), &chatpb.CreateConversationRequest{
		Name:      requestBody.Name,
		UserIds:   requestBody.UserIDs,
		CreatorId: requestBody.CreatorID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"conversation": resp,
	})
}

func (h *ChatHandler) GetMessages(c *gin.Context) {
	conversationID := c.Query("conversation_id")

	resp, err := h.Client.GetMessages(context.Background(), &chatpb.GetMessagesRequest{
		ConversationId: conversationID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": resp.Messages,
		"total_count": resp.TotalCount,
	})
}

func (h *ChatHandler) SendMessage(c *gin.Context) {
	var requestBody struct {
		ConversationID      string   	`json:"conversation_id" binding:"required"`
		SenderID   			string 		`json:"sender_id" binding:"required"`
		Content 			string   	`json:"content" binding:"required"`
		MediaType 			string		`json:"media_type" binding:"required"`
		MediaUrl 			string		`json:"media_url"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.Client.SendMessage(context.Background(), &chatpb.SendMessageRequest{
		ConversationId: requestBody.ConversationID,
		SenderId: requestBody.SenderID,
		Content: requestBody.Content,
		MediaType: requestBody.MediaType,
		MediaUrl: requestBody.MediaUrl,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": resp,
	})
}