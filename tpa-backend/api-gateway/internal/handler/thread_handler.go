package handler

import (
	"context"
	"log"
	"net/http"

	// pb "github.com/Wenfuuu/tpa-web/service-user/proto"
	threadpb "github.com/Wenfuuu/tpa-web/pkg/proto/thread"
	"github.com/gin-gonic/gin"
)

type ThreadHandler struct {
	Client threadpb.ThreadServiceClient
}

func NewThreadHandler(threadClient threadpb.ThreadServiceClient) *ThreadHandler {
	return &ThreadHandler{
		Client: threadClient,
	}
}

func (h *ThreadHandler) GetAllThreads(c *gin.Context) {
	resp, err := h.Client.GetAllThreads(context.Background(), &threadpb.GetAllThreadsRequest{})
		
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, resp.Threads)
}

func (h *ThreadHandler) PostThread(c *gin.Context) {
	var requestBody struct {
		UserID          string   `json:"userId" binding:"required"`       // UUID format
		CategoryID      string   `json:"categoryId" binding:"required"`   // UUID format
		CommunityID     string   `json:"communityId" binding:"required"`  // UUID format
		Content         string   `json:"content" binding:"required"`
		ReplyPermission string   `json:"replyPermission" binding:"required"`   // e.g. "everyone"
		IsAdvertisement bool     `json:"isAdvertisement"`
		ScheduledAt     string   `json:"scheduledAt"`                          // Optional; RFC3339 format
		ImageUrls       []string `json:"imageUrls"`                            // Array of image URLs
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println(requestBody.UserID)
	log.Println(requestBody.CategoryID)
	log.Println(requestBody.CommunityID)
	log.Println(requestBody.Content)
	resp, err := h.Client.PostThread(context.Background(), &threadpb.PostThreadRequest{
		UserId:          requestBody.UserID,
		CategoryId:      requestBody.CategoryID,
		CommunityId:     requestBody.CommunityID,
		Content:         requestBody.Content,
		ReplyPermission: requestBody.ReplyPermission,
		IsAdvertisement: requestBody.IsAdvertisement,
		ScheduledAt:     requestBody.ScheduledAt,
		ImageUrls:       requestBody.ImageUrls,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"thread": resp.Thread,
	})
}