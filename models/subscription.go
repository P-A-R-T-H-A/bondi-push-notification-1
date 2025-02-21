package models

import "time"

type PushSubscribers struct {
	Id        int       `orm:"auto;pk"`
	StudentId string    `orm:"unique"`
	Endpoint  string    `orm:"unique"`
	Auth      string    `orm:"unique"`
	P256dh    string    `orm:"unique"`
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
}
