package models

import "time"

type StudentCourse struct {
	Id                  int    `orm:"auto;pk"`
	StudentId           string `orm:"index"`
	CourseId            string `orm:"index"`
	IsBanned            bool
	IsActive            bool
	IsFbGroupJoined     string
	IsDoubtSolvedJoined bool
	PreviousCourseIds   string
	MaterialAddress     string    `orm:"type(text)"`
	CreatedAt           time.Time `orm:"auto_now_add;type(datetime)"`
	UpdatedAt           time.Time `orm:"auto_now;type(datetime)"`
}

type PushNotification struct {
	Id                  int `orm:"auto;pk"`
	CourseId            string
	CourseName          string
	CreatorId           string
	CreatorName         string
	NotificationContent string `orm:"type(text)"`
	NotificationImage   string `orm:"null"`
	Url     string `orm:"null"`
	Status              string
	CreatedAt           time.Time `orm:"auto_now_add;type(datetime)"`
	UpdatedAt           time.Time `orm:"auto_now;type(datetime)"`
}

type PushSubscribers struct {
	Id        int       `orm:"auto;pk"`
	StudentId string    `orm:"null"`
	Endpoint  string    `orm:"null"`
	Auth      string    `orm:"null"`
	P256dh    string    `orm:"null"`
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
