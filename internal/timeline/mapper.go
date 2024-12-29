package timeline

import (
	"github.com/ibanezv/littletwitter/internal/tweet"
	"github.com/ibanezv/littletwitter/pkg/db"
)

func FromDB(tl *db.DbTimeline) *TimelineUser {
	tweetsList := make([]tweet.Tweet, 0)
	for i := 0; i < len(tl.Tweets); i++ {
		t := tweet.Tweet{
			UserId: tl.Tweets[i].UserId,
			Text:   tl.Tweets[i].Text,
			Id:     tl.Tweets[i].Id}
		tweetsList = append(tweetsList, t)
	}
	return &TimelineUser{UserId: tl.UserId, Tweets: tweetsList}
}

func ToDB(tl *TimelineUser) *db.DbTimeline {
	dbTweetsList := make([]db.DbTweet, 0)
	for i := 0; i < len(tl.Tweets); i++ {
		t := db.DbTweet{
			Id:       tl.Tweets[i].Id,
			UserId:   tl.Tweets[i].UserId,
			Text:     tl.Tweets[i].Text,
			DateTime: tl.Tweets[i].Date,
		}
		dbTweetsList = append(dbTweetsList, t)
	}
	return &db.DbTimeline{UserId: tl.UserId, Tweets: dbTweetsList}
}
