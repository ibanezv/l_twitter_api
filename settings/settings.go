package settings

import (
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Settings struct {
		Tweets   `yaml:"tweets"`
		Cache    `yaml:"cache"`
		Timeline `yaml:"timeline"`
	}

	Tweets struct {
		Limit int `yaml:"charlimit"`
	}

	Cache struct {
		ExpirationTweets    int `yaml:"expiration_tweets"`
		ExpirationTimeline  int `yaml:"expiration_time_line"`
		ExpirationFollowers int `yaml:"expiration_followers"`
	}

	Timeline struct {
		MaxTweetsPerUser  uint64 `yaml:"tweets_per_user"`
		MaxTweetsTimeline int    `yaml:"max_tweets"`
	}
)

func NewAppSettings() (*Settings, error) {
	appSettings := &Settings{}
	_, err := os.Executable()
	if err != nil {
		log.Println(err)
	}

	err = cleanenv.ReadConfig("../../settings/settings.yml", appSettings)
	if err != nil {
		return nil, fmt.Errorf("loading settings error: %w", err)
	}

	return appSettings, nil
}
