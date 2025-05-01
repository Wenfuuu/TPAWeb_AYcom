package handler

import (
	"context"
	"net/http"

	followpb "github.com/Wenfuuu/tpa-web/pkg/proto/follow"
	"github.com/gin-gonic/gin"
)

type FollowHandler struct {
	Client followpb.FollowServiceClient
}

func NewFollowHandler(followClient followpb.FollowServiceClient) *FollowHandler {
	return &FollowHandler{
		Client: followClient,
	}
}

func (h *FollowHandler) GetFollower(c *gin.Context) {
	user_id := c.Query("user_id")

	resp, err := h.Client.GetFollower(context.Background(), &followpb.GetFollowerRequest{
		UserId: user_id,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.Follows)
}

func (h *FollowHandler) GetFollowing(c *gin.Context) {
	userId := c.Query("user_id")

	resp, err := h.Client.GetFollowing(context.Background(), &followpb.GetFollowingRequest{
		UserId: userId,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp.Follows)
}