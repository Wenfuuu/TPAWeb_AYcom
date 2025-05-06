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

func (h *ChatHandler) CreateConversation(c *gin.Context)  {
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