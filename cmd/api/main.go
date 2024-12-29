package main

import (
	"log"

	"github.com/ibanezv/littletwitter/config"
	"github.com/ibanezv/littletwitter/internal/app"
	"github.com/ibanezv/littletwitter/settings"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Application settings
	appSettings, err := settings.NewAppSettings()
	if err != nil {
		log.Fatalf("App settings error: %s", err)
	}

	// Run
	app.Run(cfg, appSettings)
}
