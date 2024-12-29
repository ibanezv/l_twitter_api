package tweet

import (
	"errors"

	"github.com/ibanezv/littletwitter/pkg/logger"
	"github.com/ibanezv/littletwitter/pkg/repository"
	"github.com/ibanezv/littletwitter/settings"
)

var ErrExceedingCharLimit = errors.New("exceeding message max length")

type TwitterService interface {
	Add(*Tweet) (*Tweet, error)
	Get(int64) (*Tweet, error)
}

type Twitter struct {
	repo repository.Repo
	log  logger.Interface
	app  *settings.Settings
}

func NewTwitter(repo repository.Repo, log logger.Interface, app *settings.Settings) TwitterService {
	return Twitter{repo, log, app}
}

func (t Twitter) Add(tt *Tweet) (*Tweet, error) {

	if len(tt.Text) > t.app.Tweets.Limit {
		return nil, ErrExceedingCharLimit
	}

	newtt, err := t.repo.SaveTweet(TweetToDb(tt), t.app.Timeline.MaxTweetsPerUser)
	if err != nil {
		t.log.Error(err.Error())
		return nil, err
	}

	return TweetFromDb(newtt), nil
}

func (t Twitter) Get(id int64) (*Tweet, error) {
	tt, err := t.repo.GetTweet(id)
	if err != nil {
		return nil, err
	}
	return TweetFromDb(tt), nil
}
