package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ibanezv/littletwitter/internal/follower"
)

type FollowersHandler struct {
	service follower.FollowerService
}

func NewFollowerHandler(serv follower.FollowerService) *FollowersHandler {
	return &FollowersHandler{serv}
}

type doFollowerRequest struct {
	UserId         int64     `json:"user_id" binding:"required"  example:"1"`
	UserIdFollowed int64     `json:"user_id_followed" binding:"required"  example:"2"`
	StartDate      time.Time `json:"start_date" type:"string" format:"date-time" binding:"required"  example:"2024-12-21T17:32:28Z"`
}

type response struct {
	Message string `json:"message"`
}

// @Summary     Set a new follower
// @Description Create a new follower
// @Tags  	    follower
// @Accept      json
// @Produce     json
// @Success     201
// @Failure     500 {object} response
// @Router      /follow [post]
func (h *FollowersHandler) PostFollower(ctx *gin.Context) {
	req := follower.FollowerUsers{}
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}

	f := follower.FollowerUsers{UserId: req.UserId, UserIdFollowed: req.UserIdFollowed, StartDate: req.StartDate}
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
