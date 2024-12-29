package tweet

import "time"

type Tweet struct {
	Id     int64     `json:"id,omitempty"`
	UserId int64     `json:"user_id"`
	Text   string    `json:"text"`
	Date   time.Time `json:"date_time,omitempty"`
}
