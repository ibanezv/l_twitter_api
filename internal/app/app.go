package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/ibanezv/littletwitter/config"
	"github.com/ibanezv/littletwitter/pkg/cache"
	"github.com/ibanezv/littletwitter/pkg/httpserver"
	"github.com/ibanezv/littletwitter/pkg/logger"
	postgresdb "github.com/ibanezv/littletwitter/pkg/postgres"
	redisCache "github.com/ibanezv/littletwitter/pkg/redis"
	"github.com/ibanezv/littletwitter/pkg/repository"
	"github.com/ibanezv/littletwitter/settings"
)

func Run(cfg *config.Config, appSettings *settings.Settings) {
	l := logger.New(cfg.Log.Level)

	//DB repository
	redisClient, err := redisCache.NewClient(cfg, appSettings)
	if err != nil {
		l.Error("error trying to connect to cache %s", err.Error())
	}
	client := cache.NewDbCacheEngine(redisClient)

	posgresConn, err := postgresdb.NewDbEngine(cfg)
	defer posgresConn.Close()
	if err != nil {
		l.Fatal("DB connection fail %+v", err)
	}

	repo := repository.NewRepository(posgresConn, client, l)

	// HTTP Server
	handler := gin.New()
	NewRouter(handler, l, repo, appSettings)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err := <-httpServer.Notify():
		l.Error(fmt.Errorf("littleTwiterapp - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("littleTwiterapp - httpServer.Shutdown: %w", err))
	}
}
