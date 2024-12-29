package mocks

import (
	"errors"

	"github.com/ibanezv/littletwitter/pkg/db"
	"github.com/ibanezv/littletwitter/settings"
)

type MockDBStorage struct {
}

type MockCache struct {
	CountGetTimeline int
}

func CreateMockDB() *MockDBStorage {
	return &MockDBStorage{}
}

func CreateMockCache() *MockCache {
	return &MockCache{CountGetTimeline: 0}
}

func (m *MockCache) GetTweet(int64) (*db.DbTweet, error) {
	return &db.DbTweet{Id: 1, UserId: 1}, nil
}

func (m *MockCache) SetTweet(*db.DbTweet) error {
	return nil
}

func (m *MockCache) GetTimeline(userId int64) (*db.DbTimeline, error) {
	if userId == 1 {
		m.CountGetTimeline++
		return &db.DbTimeline{UserId: userId}, nil
	}
	return nil, nil
}

func (m *MockCache) SetTimeline(int64, *db.DbTimeline) error {
	return nil
}

func (m *MockCache) DeleteTimeline(int64) error {
	return nil
}

func (m *MockCache) SetFollowers(int64, []int64) error {
	return nil
}

func (m *MockCache) GetFollowers(int64) ([]int64, error) {
	return []int64{2}, nil
}

func (m *MockDBStorage) GetTweet(int64) (*db.DbTweet, error) {
	return &db.DbTweet{Id: 1, UserId: 1}, nil
}

func (m *MockDBStorage) SetTweet(t *db.DbTweet) (*db.DbTweet, error) {
	if t.UserId == 1 {
		return &db.DbTweet{Id: 1, UserId: t.UserId, Text: t.Text}, nil
	}
	return nil, errors.New("test error")
}

func (m *MockDBStorage) GetTweetsByUser(int64, uint64) ([]db.DbTweet, error) {
	return []db.DbTweet{{Id: 1, UserId: 1}}, nil
}

func (m *MockDBStorage) SetFollower(f *db.DbFollower) error {
	if f.UserId == 3 && f.UserFollowedId == 4 {
		return errors.New("test")
	}
	return nil
}

func (m *MockDBStorage) GetFollower(f *db.DbFollower) ([]db.DbFollower, error) {
	if f.UserId == 3 && f.UserFollowedId == 4 {
		return []db.DbFollower{{UserId: f.UserId, UserFollowedId: f.UserFollowedId}}, nil
	}
	if f.UserId == 5 && f.UserFollowedId == 6 {
		return nil, errors.New("test")
	}
	return nil, nil
}

func (m *MockDBStorage) GetFollowers(int64) ([]int64, error) {
	return []int64{1, 3}, nil
}

func (m *MockDBStorage) GetFollowing(userIdFollower int64) ([]int64, error) {
	if userIdFollower == 2 {
		return nil, errors.New("testing error")
	}
	return []int64{2}, nil
}

func (m *MockDBStorage) createTweetsTimeLine(userId int64, _ uint64) (*db.DbTimeline, error) {
	if userId == 2 {
		return nil, errors.New("test error")
	}
	return &db.DbTimeline{UserId: userId, Tweets: []db.DbTweet{{UserId: userId, Text: "text test"}}}, nil
}

func CreateSettingsMock() *settings.Settings {
	return &settings.Settings{
		Tweets:   settings.Tweets{Limit: 10},
		Cache:    settings.Cache{ExpirationTweets: 1, ExpirationTimeline: 1, ExpirationFollowers: 1},
		Timeline: settings.Timeline{MaxTweetsPerUser: 10, MaxTweetsTimeline: 20},
	}
}
