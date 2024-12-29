package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ibanezv/littletwitter/internal/tweet"
)

type TwitterHandler struct {
	service tweet.TwitterService
}

func NewTwitterHander(serv tweet.TwitterService) *TwitterHandler {
	return &TwitterHandler{serv}
}

func (h *TwitterHandler) PostTwitter(ctx *gin.Context) {
	bodyReq := tweet.Tweet{}
	err := ctx.BindJSON(&bodyReq)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
	}

	tt, err := h.service.Add(&bodyReq)
	if err != nil {
		if errors.Is(err, tweet.ErrExceedingCharLimit) {
			ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "message exceed max length"})
			return
		}
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	ctx.IndentedJSON(http.StatusCreated, tt)
}
