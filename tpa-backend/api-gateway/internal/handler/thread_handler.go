package handler

import (
	"context"
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

func(h *ThreadHandler) GetAllThreads(c *gin.Context) {
	resp, err := h.Client.GetAllThreads(context.Background(), &threadpb.GetAllThreadsRequest{})
		
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, resp.Threads)
}