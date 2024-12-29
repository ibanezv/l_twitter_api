package db

import (
	"time"
)

type DbTweet struct {
	Id       int64     `json:"id"`
	UserId   int64     `json:"user_id"`
	Text     string    `json:"text"`
	DateTime time.Time `json:"datetime"`
}

type DbFollower struct {
	Id             int64     `json:"id"`
	UserId         int64     `json:"user_id"`
	UserFollowedId int64     `json:"user_id_followed"`
	DateInit       time.Time `json:"datetime"`
}

type DbTimeline struct {
	UserId int64     `json:"user_id"`
	Tweets []DbTweet `json:"tweets"`
}
