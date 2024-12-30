package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ibanezv/littletwitter/internal/tweet"
)

type TwitterHandler struct {
	service tweet.TwitterService
}

func NewTwitterHander(serv tweet.TwitterService) *TwitterHandler {
	return &TwitterHandler{serv}
}

type doTweetRequest struct {
	UserId int64     `json:"user_id" binding:"required"  example:"1"`
	Text   string    `json:"text" binding:"required"  example:"message example"`
	Date   time.Time `json:"date_time" type:"string" format:"date-time" binding:"required"  example:"2017-07-21T17:32:28Z"`
}

type tweetResponse struct {
	Id     int64     `json:"id"`
	UserId int64     `json:"user_id"`
	Text   string    `json:"text"`
	Date   time.Time `json:"date_time"`
}

// @Summary     New Tweet
// @Description Create a new tweet
// @Tags  	    tweet
// @Accept      json
// @Produce     json
// @Param       request body doTweetRequest true "New tweet"
// @Success     201 {object} tweetResponse
// @Failure     500 {object} response
// @Router      /tweet [post]
func (h *TwitterHandler) PostTwitter(ctx *gin.Context) {
	bodyReq := doTweetRequest{}
	err := ctx.BindJSON(&bodyReq)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}

	tweetReq := tweet.Tweet{UserId: bodyReq.UserId, Text: bodyReq.Text, Date: bodyReq.Date}
	tt, err := h.service.Add(&tweetReq)
	if err != nil {
		if errors.Is(err, tweet.ErrExceedingCharLimit) {
			ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "message exceed max length"})
			return
		}
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	ctx.IndentedJSON(http.StatusCreated, tweetResponse{Id: tt.Id, UserId: tt.UserId, Text: tt.Text, Date: tt.Date})
}
