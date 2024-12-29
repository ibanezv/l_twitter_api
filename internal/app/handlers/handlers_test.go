package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ibanezv/littletwitter/internal/follower"
	"github.com/ibanezv/littletwitter/internal/mocks"
	"github.com/ibanezv/littletwitter/internal/timeline"
	"github.com/ibanezv/littletwitter/internal/tweet"
	"github.com/ibanezv/littletwitter/pkg/cache"
	"github.com/ibanezv/littletwitter/pkg/db"
	"github.com/ibanezv/littletwitter/pkg/logger"
	"github.com/ibanezv/littletwitter/pkg/repository"
	"github.com/stretchr/testify/assert"
)

func TestFollowersHandler(t *testing.T) {
	var tests = []struct {
		name           string
		follower       follower.FollowerUsers
		mockCache      cache.Cache
		mockDB         db.DBStorage
		statusExpected int
	}{
		{
			name:           "add sucess",
			follower:       follower.FollowerUsers{UserId: 1, UserIdFollowed: 2},
			mockCache:      mocks.CreateMockCache(),
			mockDB:         mocks.CreateMockDB(),
			statusExpected: http.StatusCreated,
		},
		{
			name:           "fail database error",
			follower:       follower.FollowerUsers{UserId: 5, UserIdFollowed: 6},
			mockCache:      mocks.CreateMockCache(),
			mockDB:         mocks.CreateMockDB(),
			statusExpected: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			body, _ := json.Marshal(test.follower)
			ctx.Request = &http.Request{
				Body: io.NopCloser(bytes.NewBuffer(body)),
			}

			l := logger.New("test")
			repo := repository.NewRepository(test.mockDB, test.mockCache, l)
			serv := follower.NewFollowerService(repo, l)
			follHandler := NewFollowerHandler(serv)

			// when
			follHandler.PostFollower(ctx)

			//then
			assert.Equal(t, test.statusExpected, ctx.Writer.Status())
		})
	}
}

func TestTweetsHandler(t *testing.T) {
	var tests = []struct {
		name           string
		tweet          tweet.Tweet
		mockCache      cache.Cache
		mockDB         db.DBStorage
		statusExpected int
	}{
		{
			name:           "add sucess",
			tweet:          tweet.Tweet{UserId: 1, Text: "test"},
			mockCache:      mocks.CreateMockCache(),
			mockDB:         mocks.CreateMockDB(),
			statusExpected: http.StatusCreated,
		},
		{
			name:           "fail database error",
			tweet:          tweet.Tweet{UserId: 2, Text: "test"},
			mockCache:      mocks.CreateMockCache(),
			mockDB:         mocks.CreateMockDB(),
			statusExpected: http.StatusInternalServerError,
		},
		{
			name:           "fail text message exceed char limit",
			tweet:          tweet.Tweet{UserId: 1, Text: "test message too long"},
			mockCache:      mocks.CreateMockCache(),
			mockDB:         mocks.CreateMockDB(),
			statusExpected: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			body, _ := json.Marshal(test.tweet)
			ctx.Request = &http.Request{
				Body: io.NopCloser(bytes.NewBuffer(body)),
			}

			l := logger.New("test")
			repo := repository.NewRepository(test.mockDB, test.mockCache, l)
			appSettings := mocks.CreateSettingsMock()
			serv := tweet.NewTwitter(repo, l, appSettings)
			tweetHandler := NewTwitterHander(serv)

			// when
			tweetHandler.PostTwitter(ctx)

			//then
			assert.Equal(t, test.statusExpected, ctx.Writer.Status())
		})
	}
}

func TestTimelineHandler(t *testing.T) {
	var tests = []struct {
		name           string
		userId         string
		mockCache      cache.Cache
		mockDB         db.DBStorage
		statusExpected int
	}{
		{
			name:           "add sucess",
			userId:         "1",
			mockCache:      mocks.CreateMockCache(),
			mockDB:         mocks.CreateMockDB(),
			statusExpected: http.StatusOK,
		},
		{
			name:           "fail database error",
			userId:         "2",
			mockCache:      mocks.CreateMockCache(),
			mockDB:         mocks.CreateMockDB(),
			statusExpected: http.StatusInternalServerError,
		},
		{
			name:           "fail bad request",
			userId:         "2aa",
			mockCache:      mocks.CreateMockCache(),
			mockDB:         mocks.CreateMockDB(),
			statusExpected: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Params = append(ctx.Params, gin.Param{Key: "userId", Value: test.userId})

			l := logger.New("test")
			repo := repository.NewRepository(test.mockDB, test.mockCache, l)
			appSettings := mocks.CreateSettingsMock()
			serv := timeline.NewTimeLine(repo, l, appSettings)
			timelineHandler := NewTimelineHandler(serv)

			// when
			timelineHandler.GetTimeline(ctx)

			//then
			assert.Equal(t, test.statusExpected, ctx.Writer.Status())
		})
	}
}
