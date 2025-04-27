package handler

import (
	"context"
	"net/http"

	// pb "github.com/Wenfuuu/tpa-web/service-user/proto"
	userpb "github.com/Wenfuuu/tpa-web/pkg/proto/user"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Client userpb.UserServiceClient
}

func NewUserHandler(userClient userpb.UserServiceClient) *UserHandler {
	return &UserHandler{
		Client: userClient,
	}
}

func(h *UserHandler) GetAllUsers(c *gin.Context) {
	resp, err := h.Client.GetAllUsers(context.Background(), &userpb.GetAllUsersRequest{})
		
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, resp.Users)
}