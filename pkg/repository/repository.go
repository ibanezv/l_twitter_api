package repository

import (
	"sort"

	"github.com/ibanezv/littletwitter/pkg/cache"
	"github.com/ibanezv/littletwitter/pkg/db"
	"github.com/ibanezv/littletwitter/pkg/logger"
)

type Repo interface {
	GetTweet(int64) (*db.DbTweet, error)
	SaveTweet(*db.DbTweet, uint64) (*db.DbTweet, error)
	GetTimeline(int64, uint64) (*db.DbTimeline, error)
	SetFollower(*db.DbFollower) error
	GetFollower(*db.DbFollower) ([]db.DbFollower, error)
}

type Repository struct {
	DB     db.DBStorage
	Cache  cache.Cache
	Logger logger.Interface
}

func NewRepository(db db.DBStorage, cache cache.Cache, logger logger.Interface) *Repository {
	return &Repository{db, cache, logger}
}

// GetTweet() returns a tweet by id
func (r *Repository) GetTweet(id int64) (*db.DbTweet, error) {
	tweet, err := r.Cache.GetTweet(id)
	if err == nil {
		if tweet != nil {
			return tweet, err
		}
	} else {
		r.Logger.Error("redis error: %s", err.Error())
	}
	return r.DB.GetTweet(id)
}

// SaveTweet() save a new tweet and refresh timelines cached for each user (follower) afected
func (r *Repository) SaveTweet(tweet *db.DbTweet, topByUser uint64) (*db.DbTweet, error) {
	newTweet, err := r.DB.SetTweet(tweet)
	if err != nil {
		return nil, err
	}
	go r.Cache.SetTweet(newTweet)
	go r.refreshTimelines(tweet.UserId, topByUser)
	return newTweet, nil
}

// GetTimeline() returns timeline for a user
func (r *Repository) GetTimeline(userId int64, topByUser uint64) (*db.DbTimeline, error) {
	dbTimeline, err := r.Cache.GetTimeline(userId)
	if err == nil {
		if dbTimeline != nil {
			return dbTimeline, nil
		}
	} else {
		r.Logger.Error("redis error: %s", err.Error())
	}

	tl, err := r.createTweetsTimeLine(userId, topByUser)
	if err != nil {
		return nil, err
	}

	r.Cache.SetTimeline(userId, tl)
	return tl, err
}

// SetFollower() add new follower to db
func (r *Repository) SetFollower(follower *db.DbFollower) error {
	err := r.DB.SetFollower(follower)
	if err != nil {
		return err
	}

	go r.refreshFollowers(follower.UserFollowedId)
	return nil
}

func (r *Repository) GetFollower(follower *db.DbFollower) ([]db.DbFollower, error) {
	return r.DB.GetFollower(follower)
}

// refreshFollowers() refresh followers list from cache
func (r *Repository) refreshFollowers(userIdFollowed int64) error {
	followers, err := r.DB.GetFollowers(userIdFollowed)
	if err != nil {
		return err
	}

	err = r.Cache.SetFollowers(userIdFollowed, followers)
	if err != nil {
		return err
	}

	return nil
}

// refreshTimelines() refresh timelines from cache for each follower afected
func (r *Repository) refreshTimelines(userIdFollowed int64, topByUser uint64) error {
	followers, err := r.getFollowers(userIdFollowed)
	if err != nil {
		return err
	}

	for i := 0; i < len(followers); i++ {
		err = r.Cache.DeleteTimeline(followers[i])
		tl, _ := r.createTweetsTimeLine(followers[i], topByUser)
		r.Cache.SetTimeline(followers[i], tl)
	}

	return nil
}

// createTweetsTimeLine() returns timeline from Database for a user
func (r *Repository) createTweetsTimeLine(userId int64, topByUser uint64) (*db.DbTimeline, error) {
	following, err := r.DB.GetFollowing(userId)
	if err != nil {
		return nil, err
	}

	tweets := make([]db.DbTweet, 0)
	for i := 0; i < len(following); i++ {
		t, err := r.DB.GetTweetsByUser(following[i], topByUser)
		if err != nil {
			return nil, err
		}
		tweets = append(tweets, t...)
	}

	sort.Slice(tweets, func(i, j int) bool { return tweets[i].Id > tweets[j].Id })

	return &db.DbTimeline{UserId: userId, Tweets: tweets}, nil
}

// getFollowers() get followers for a specific user
func (r *Repository) getFollowers(userIdFollowed int64) ([]int64, error) {
	followers, err := r.Cache.GetFollowers(userIdFollowed)
	if err == nil {
		if followers != nil {
			return followers, nil
		}
	} else {
		r.Logger.Error("redis error: %s", err.Error())
	}

	return r.DB.GetFollowers(userIdFollowed)
}
