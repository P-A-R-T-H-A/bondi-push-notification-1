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

type PushNotificationSubscription struct {
	StudentId    string
	Notification NotificationSubscription
	Success      bool
	Error        string
	Message      string
}

type NotificationSubscription struct {
	Endpoint       string
	ExpirationTime any
	Keys           Keys
}
type Keys struct {
	P256Dh string
	Auth   string
}
