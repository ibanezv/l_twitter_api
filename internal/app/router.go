package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/ibanezv/littletwitter/docs"
	"github.com/ibanezv/littletwitter/internal/app/handlers"
	"github.com/ibanezv/littletwitter/internal/follower"
	"github.com/ibanezv/littletwitter/internal/timeline"
	"github.com/ibanezv/littletwitter/internal/tweet"
	"github.com/ibanezv/littletwitter/pkg/logger"
	"github.com/ibanezv/littletwitter/pkg/repository"
	"github.com/ibanezv/littletwitter/settings"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter -.
// Swagger spec:
// @title
// @description
// @version     1.0
// @host
// @BasePath    /v1
func NewRouter(handler *gin.Engine, l logger.Interface, repo repository.Repo, appSettings *settings.Settings) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	handler.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
		c.IndentedJSON(http.StatusOK, gin.H{"message": "health ok"})
	})

	// Routers
	followService := follower.NewFollowerService(repo, l)
	followerHandler := handlers.NewFollowerHandler(followService)

	twitterService := tweet.NewTwitter(repo, l, appSettings)
	twitterHandler := handlers.NewTwitterHander(twitterService)

	timelineService := timeline.NewTimeLine(repo, l, appSettings)
	timelineHandler := handlers.NewTimelineHandler(timelineService)
	h := handler.Group("/v1/api-twitter")
	{
		h.POST("/follow", followerHandler.PostFollower)
		h.POST("/tweet", twitterHandler.PostTwitter)
		h.GET("/timeline/:userId", timelineHandler.GetTimeline)
	}
}
