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

func (h *ThreadHandler) GetReplies(c *gin.Context) {
	pageStr := c.Query("page")
	threadId := c.Query("thread_id")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	resp, err := h.Client.GetReplies(context.Background(), &threadpb.GetRepliesRequest{
		Page: int32(page),
		ThreadId: threadId,
	})
		
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, resp.Threads)
}

func (h *ThreadHandler) GetThreadById(c *gin.Context) {
	threadId := c.Query("thread_id")

	resp, err := h.Client.GetThreadById(context.Background(), &threadpb.GetThreadByIdRequest{
		ThreadId: threadId,
	})
		
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, resp.Thread)
}

func (h *ThreadHandler) GetTopHashtags(c *gin.Context) {
	resp, err := h.Client.GetTopHashtags(context.Background(), &threadpb.GetTopHashtagsRequest{})
		
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, resp.HashtagCount)
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

func (h *ThreadHandler) GetCommunityThreads(c *gin.Context) {
	pageStr := c.Query("page")
	communityId := c.Query("community_id")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	resp, err := h.Client.GetCommunityThreads(context.Background(), &threadpb.GetCommunityThreadsRequest{
		Page: int32(page),
		CommunityId: communityId,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.Threads)
}

func (h *ThreadHandler) GetFollowingThreads(c *gin.Context) {
	pageStr := c.Query("page")
	userId := c.Query("user_id")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	resp, err := h.Client.GetFollowingThreads(context.Background(), &threadpb.GetFollowingThreadsRequest{
		Page: int32(page),
		UserId: userId,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.Threads)
}

func (h *ThreadHandler) GetBookmark(c *gin.Context) {
	pageStr := c.Query("page")
	userId := c.Query("user_id")
	search := c.Query("search")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	resp, err := h.Client.GetBookmark(context.Background(), &threadpb.GetBookmarkRequest{
		Page: int32(page),
		UserId: userId,
		Search: search,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.Threads)
}

func (h *ThreadHandler) GetRepostThreads(c *gin.Context) {
	userId := c.Query("user_id")

	resp, err := h.Client.GetRepostThreads(context.Background(), &threadpb.GetRepostThreadsRequest{
		UserId: userId,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.Threads)
}

func (h *ThreadHandler) GetLikedThreads(c *gin.Context) {
	userId := c.Query("user_id")

	resp, err := h.Client.GetLikedThreads(context.Background(), &threadpb.GetLikedThreadsRequest{
		UserId: userId,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.Threads)
}

func (h *ThreadHandler) GetBookmarkThreads(c *gin.Context) {
	userId := c.Query("user_id")

	resp, err := h.Client.GetBookmarkThreads(context.Background(), &threadpb.GetBookmarkThreadsRequest{
		UserId: userId,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.Threads)
}

func (h *ThreadHandler) UnrepostThread(c *gin.Context) {
	var requestBody struct {
		UserID          string   `json:"userId" binding:"required"`
		ThreadID      	string   `json:"threadId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.Client.UnrepostThread(context.Background(), &threadpb.UnrepostThreadRequest{
		UserId: requestBody.UserID,
		ThreadId: requestBody.ThreadID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": resp.Message,
	})
}

func (h *ThreadHandler) RepostThread(c *gin.Context) {
	var requestBody struct {
		UserID          string   `json:"userId" binding:"required"`
		RepostOfID      string   `json:"repostOfId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.Client.RepostThread(context.Background(), &threadpb.RepostThreadRequest{
		UserId: requestBody.UserID,
		ThreadId: requestBody.RepostOfID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"thread": resp.Thread,
	})
}

func (h *ThreadHandler) ReplyThread(c *gin.Context) {
	var requestBody struct {
		UserID          string   `json:"userId" binding:"required"`       // UUID format
		ReplyToID      string   `json:"replyToId" binding:"required"`
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

	resp, err := h.Client.ReplyThread(context.Background(), &threadpb.ReplyThreadRequest{
		UserId: requestBody.UserID,
		ReplyToId: requestBody.ReplyToID,
		Content: requestBody.Content,
		ReplyPermission: requestBody.ReplyPermission,
		IsAdvertisement: requestBody.IsAdvertisement,
		ScheduledAt: requestBody.ScheduledAt,
		ImageUrls: requestBody.ImageUrls,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"thread": resp.Thread,
	})
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

func (h *ThreadHandler) LikeThread(c *gin.Context) {
	var requestBody struct {
		UserID          string   `json:"userId" binding:"required"`       // UUID format
		ThreadID      	string   `json:"threadId" binding:"required"`   // UUID format
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.Client.LikeThread(context.Background(), &threadpb.LikeThreadRequest{
		UserId: requestBody.UserID,
		ThreadId: requestBody.ThreadID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": resp.Message,
	})
}

func (h *ThreadHandler) UnlikeThread(c *gin.Context) {
	var requestBody struct {
		UserID          string   `json:"userId" binding:"required"`       // UUID format
		ThreadID      	string   `json:"threadId" binding:"required"`   // UUID format
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.Client.UnlikeThread(context.Background(), &threadpb.UnlikeThreadRequest{
		UserId: requestBody.UserID,
		ThreadId: requestBody.ThreadID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": resp.Message,
	})
}

func (h *ThreadHandler) BookmarkThread(c *gin.Context) {
	var requestBody struct {
		UserID          string   `json:"userId" binding:"required"`
		ThreadID      	string   `json:"threadId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.Client.BookmarkThread(context.Background(), &threadpb.BookmarkThreadRequest{
		UserId: requestBody.UserID,
		ThreadId: requestBody.ThreadID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": resp.Message,
	})
}

func (h *ThreadHandler) UnbookmarkThread(c *gin.Context) {
	var requestBody struct {
		UserID          string   `json:"userId" binding:"required"`
		ThreadID      	string   `json:"threadId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.Client.UnbookmarkThread(context.Background(), &threadpb.UnbookmarkThreadRequest{
		UserId: requestBody.UserID,
		ThreadId: requestBody.ThreadID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": resp.Message,
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

func (h *ThreadHandler) GetCommunityById(c *gin.Context) {
	communityId := c.Query("community_id")

	resp, err := h.Client.GetCommunityById(context.Background(), &threadpb.GetCommunityByIdRequest{
		CommunityId: communityId,
	})
		
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, resp.Community)
}

func (h *ThreadHandler) GetAllCommunity(c *gin.Context) {
	resp, err := h.Client.GetAllCommunity(context.Background(), &threadpb.GetAllCommunitiesRequest{})
		
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, resp.Communities)
}

func (h *ThreadHandler) GetPendingCommunity(c *gin.Context) {
	userId := c.Query("user_id")
	pageStr := c.Query("page")
	limitStr := c.Query("limit")
	search := c.Query("search")
	category_id := c.Query("category_id")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 25
	}

	resp, err := h.Client.GetPendingCommunity(context.Background(), &threadpb.GetPendingCommunityRequest{
		UserId: userId,
		Page: int32(page),
		Limit: int32(limit),
		Search: search,
		CategoryId: category_id,
	})
		
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, resp.Communities)
}

func (h *ThreadHandler) GetJoinedCommunity(c *gin.Context) {
	userId := c.Query("user_id")
	pageStr := c.Query("page")
	limitStr := c.Query("limit")
	search := c.Query("search")
	category_id := c.Query("category_id")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 25
	}

	resp, err := h.Client.GetJoinedCommunity(context.Background(), &threadpb.GetJoinedCommunityRequest{
		UserId: userId,
		Page: int32(page),
		Limit: int32(limit),
		Search: search,
		CategoryId: category_id,
	})
		
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, resp.Communities)
}

func (h *ThreadHandler) GetCommunityPagination(c *gin.Context) {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")
	search := c.Query("search")
	category_id := c.Query("category_id")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 25
	}

	resp, err := h.Client.GetCommunityPagination(context.Background(), &threadpb.GetCommunityPaginationRequest{
		Page: int32(page),
		Limit: int32(limit),
		Search: search,
		CategoryId: category_id,
	})
		
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, resp.Communities)
}

func (h *ThreadHandler) JoinCommunity(c *gin.Context) {
	var requestBody struct {
		UserID          string   	`json:"userId" binding:"required"`
		CommunityID     string   	`json:"communityId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.Client.JoinCommunity(context.Background(), &threadpb.JoinCommunityRequest{
		UserId: requestBody.UserID,
		CommunityId: requestBody.CommunityID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": resp.Message,
	})
}

func (h *ThreadHandler) RequestJoinCommunity(c *gin.Context) {
	var requestBody struct {
		UserID          string   	`json:"userId" binding:"required"`
		CommunityID     string   	`json:"communityId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.Client.RequestJoinCommunity(context.Background(), &threadpb.RequestJoinCommunityRequest{
		UserId: requestBody.UserID,
		CommunityId: requestBody.CommunityID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": resp.Message,
	})
}

func (h *ThreadHandler) GetAllCommunityCategories(c *gin.Context) {
	resp, err := h.Client.GetAllCommunityCategories(context.Background(), &threadpb.GetAllCommunityCategoriesRequest{})
		
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, resp.CommunityCategories)
}

func (h *ThreadHandler) GetAllThreadFile(c *gin.Context) {
	resp, err := h.Client.GetAllThreadFile(context.Background(), &threadpb.GetAllThreadFileRequest{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, resp.ThreadFiles)
}