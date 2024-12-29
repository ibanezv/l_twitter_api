package redisCache

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/ibanezv/littletwitter/config"
	"github.com/ibanezv/littletwitter/pkg/db"
	"github.com/ibanezv/littletwitter/settings"
)

const (
	redisTimelinePrefix = "tl-"
	redisTwitterPrefix  = "tt-"
	redisFollowerPrefix = "fw-"
)

type RedisClient struct {
	client   *redis.Client
	settings *settings.Settings
}

func NewClient(cfg *config.Config, appSettings *settings.Settings) (*RedisClient, error) {
	c := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Address,
		DB:   cfg.Redis.Db,
	})

	_, err := c.Ping().Result()

	return &RedisClient{client: c, settings: appSettings}, err
}

func (c *RedisClient) GetTweet(id int64) (*db.DbTweet, error) {
	var t db.DbTweet
	cmd := c.client.Get(redisTwitterPrefix + strconv.FormatInt(id, 10))
	if len(cmd.Val()) == 0 {
		return nil, nil
	}

	err := json.Unmarshal([]byte(cmd.Val()), &t)
	return &t, err
}

func (c *RedisClient) SetTweet(tweet *db.DbTweet) error {
	data, err := json.Marshal(*tweet)
	if err != nil {
		return errors.New("error marshaling tweet: " + err.Error())
	}

	return c.client.Set(redisTwitterPrefix+strconv.FormatInt(tweet.Id, 10), data, time.Duration(c.settings.ExpirationTweets)*time.Minute).Err()
}

func (c *RedisClient) GetTimeline(userId int64) (*db.DbTimeline, error) {
	var tl db.DbTimeline
	cmd := c.client.Get(redisTimelinePrefix + strconv.FormatInt(userId, 10))
	if len(cmd.Val()) == 0 {
		return nil, nil
	}

	err := json.Unmarshal([]byte(cmd.Val()), &tl)
	if err != nil {
		return nil, err
	}

	return &tl, nil
}

func (c *RedisClient) SetTimeline(userId int64, tl *db.DbTimeline) error {
	data, err := json.Marshal(tl)
	if err != nil {
		return errors.New("error marshaling timeline: " + err.Error())
	}
	return c.client.Set(redisTimelinePrefix+strconv.FormatInt(userId, 10), data, time.Duration(c.settings.ExpirationTimeline)*time.Minute).Err()
}

func (c *RedisClient) DeleteTimeline(userId int64) error {
	return c.client.Del(redisTimelinePrefix + strconv.FormatInt(userId, 10)).Err()
}

func (c *RedisClient) SetFollowers(userId int64, followers []int64) error {
	data, err := json.Marshal(followers)
	if err != nil {
		return errors.New("error marshaling follower list: " + err.Error())
	}
	return c.client.Set(redisFollowerPrefix+strconv.FormatInt(userId, 10), data, time.Duration(c.settings.ExpirationFollowers)*time.Minute).Err()
}

func (c *RedisClient) GetFollowers(userId int64) ([]int64, error) {
	var followers []int64
	cmd := c.client.Get(redisFollowerPrefix + strconv.FormatInt(userId, 10))
	if len(cmd.Val()) == 0 {
		return nil, nil
	}
	err := json.Unmarshal([]byte(cmd.Val()), &followers)
	return followers, err
}
