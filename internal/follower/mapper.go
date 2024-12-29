package follower

import "github.com/ibanezv/littletwitter/pkg/db"

func FollowerToDB(f *FollowerUsers) *db.DbFollower {
	return &db.DbFollower{UserId: f.UserId, UserFollowedId: f.UserIdFollowed, DateInit: f.StartDate}
}

func FollowerFromDB(f *db.DbFollower) *FollowerUsers {
	return &FollowerUsers{UserId: f.UserId, UserIdFollowed: f.UserFollowedId, StartDate: f.DateInit}
}
