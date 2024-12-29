package follower

import "time"

type FollowerUsers struct {
	UserId         int64     `json:"user_id"`
	UserIdFollowed int64     `json:"user_id_followed"`
	StartDate      time.Time `json:"start_date,omitempty"`
}
