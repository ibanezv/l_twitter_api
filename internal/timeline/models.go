package timeline

import "github.com/ibanezv/littletwitter/internal/tweet"

type TimelineUser struct {
	UserId int64         `json:"user_id"`
	Tweets []tweet.Tweet `json:"tweets"`
}
