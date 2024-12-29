package tweet

import "github.com/ibanezv/littletwitter/pkg/db"

func TweetFromDb(t *db.DbTweet) *Tweet {
	return &Tweet{Id: t.Id, UserId: t.UserId, Text: t.Text, Date: t.DateTime}
}

func TweetToDb(t *Tweet) *db.DbTweet {
	return &db.DbTweet{Id: t.Id, UserId: t.UserId, Text: t.Text, DateTime: t.Date}
}
