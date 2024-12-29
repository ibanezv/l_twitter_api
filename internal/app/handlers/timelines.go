package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ibanezv/littletwitter/internal/timeline"
)

type TimelineHandler struct {
	service timeline.TimelineService
}

func NewTimelineHandler(serv timeline.TimelineService) *TimelineHandler {
	return &TimelineHandler{serv}
}

func (h *TimelineHandler) GetTimeline(ctx *gin.Context) {
	pUserId, exists := ctx.Params.Get("userId")
	if !exists {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "param not found"})
	}

	userId, err := strconv.ParseInt(pUserId, 10, 64)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid param"})
	}

	tl, err := h.service.Get(userId)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, tl)
}
