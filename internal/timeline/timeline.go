package timeline

import (
	"github.com/ibanezv/littletwitter/pkg/logger"
	"github.com/ibanezv/littletwitter/pkg/repository"
	"github.com/ibanezv/littletwitter/settings"
)

type TimelineService interface {
	Get(int64) (*TimelineUser, error)
}

type TimeLiner struct {
	repo        repository.Repo
	log         logger.Interface
	appSettings *settings.Settings
}

func NewTimeLine(repo repository.Repo, l logger.Interface, appSettings *settings.Settings) TimelineService {
	return &TimeLiner{repo, l, appSettings}
}

func (t *TimeLiner) Get(userId int64) (*TimelineUser, error) {
	tl, err := t.repo.GetTimeline(userId, t.appSettings.Timeline.MaxTweetsPerUser)
	if err != nil {
		t.log.Error("error geting timeline %+v", err)
		return nil, err
	}

	if len(tl.Tweets) > t.appSettings.Timeline.MaxTweetsTimeline {
		tl.Tweets = tl.Tweets[:t.appSettings.Timeline.MaxTweetsTimeline]
	}

	return FromDB(tl), err
}
