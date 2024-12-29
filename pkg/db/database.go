package db

type DBStorage interface {
	GetTweet(int64) (*DbTweet, error)
	SetTweet(*DbTweet) (*DbTweet, error)
	GetTweetsByUser(int64, uint64) ([]DbTweet, error)
	SetFollower(*DbFollower) error
	GetFollower(*DbFollower) ([]DbFollower, error)
	GetFollowers(int64) ([]int64, error)
	GetFollowing(userIdFollower int64) ([]int64, error)
}

type Database struct {
	db DBStorage
}

func NewdDataBaseEngine(db DBStorage) *Database {
	return &Database{db}
}
