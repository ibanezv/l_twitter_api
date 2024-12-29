package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ibanezv/littletwitter/internal/follower"
)

type FollowersHandler struct {
	service follower.FollowerService
}

func NewFollowerHandler(serv follower.FollowerService) *FollowersHandler {
	return &FollowersHandler{serv}
}

func (h *FollowersHandler) PostFollower(ctx *gin.Context) {
	f := follower.FollowerUsers{}
	err := ctx.BindJSON(&f)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}

	err = h.service.Set(&f)
	if err != nil {
		if errors.Is(err, follower.ErrFollowerAlreadyExists) {
			ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "follower already exists"})
			return
		}
		if errors.Is(err, follower.ErrFollowerNotValid) {
			ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "follower not valid"})
		}
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	ctx.IndentedJSON(http.StatusCreated, gin.H{"message": "ok"})
}
