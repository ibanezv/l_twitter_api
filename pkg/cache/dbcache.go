package cache

import (
	"github.com/ibanezv/littletwitter/pkg/db"
)

type Cache interface {
	GetTweet(int64) (*db.DbTweet, error)
	SetTweet(*db.DbTweet) error
	GetTimeline(int64) (*db.DbTimeline, error)
	SetTimeline(int64, *db.DbTimeline) error
	DeleteTimeline(int64) error
	SetFollowers(int64, []int64) error
	GetFollowers(int64) ([]int64, error)
}

type DBCacheEngine struct {
	Cache
}

func NewDbCacheEngine(cache Cache) DBCacheEngine {
	return DBCacheEngine{Cache: cache}
}
