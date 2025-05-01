package handler

import (
	"context"
	"net/http"
	"strconv"

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
	pageStr := c.Query("page")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	resp, err := h.Client.GetAllThreads(context.Background(), &threadpb.GetAllThreadsRequest{
		Page: int32(page),
	})
		
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
		CommunityID     string   `json:"communityId"`  // UUID format
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

func (h *ThreadHandler) GetAllThreadCategories(c *gin.Context) {
	resp, err := h.Client.GetAllThreadCategories(context.Background(), &threadpb.GetAllThreadCategoriesRequest{})
		
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, resp.ThreadCategories)
}

func (h *ThreadHandler) GetAllCommunity(c *gin.Context) {
	resp, err := h.Client.GetAllCommunity(context.Background(), &threadpb.GetAllCommunitiesRequest{})
		
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, resp.Communities)
}

func (h *ThreadHandler) GetAllThreadFile(c *gin.Context) {
	resp, err := h.Client.GetAllThreadFile(context.Background(), &threadpb.GetAllThreadFileRequest{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, resp.ThreadFiles)
}