package domain

import "time"

type User struct {
	Id               int
	FirstName        string
	LastName         string
	Username         string
	ChatId           int64
	NotificationTime time.Time
}

func (u *User) NextNotificationTime() {
	u.NotificationTime = u.NotificationTime.Add(24 * time.Hour)
}
