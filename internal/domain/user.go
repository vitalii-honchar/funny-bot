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

const notificationHour = 18

func NewUser(fn string, ln string, u string, ci int64, now time.Time) *User {
	return &User{
		FirstName:        fn,
		LastName:         ln,
		Username:         u,
		ChatId:           ci,
		NotificationTime: nextNotificationTime(now),
	}
}

func nextNotificationTime(now time.Time) time.Time {
	nt := time.Date(now.Year(), now.Month(), now.Day(), notificationHour, 0, 0, 0, now.Location())

	if nt.After(now) {
		return nt
	}
	return now
}

func (u *User) NextNotificationTime() {
	u.NotificationTime = u.NotificationTime.Add(24 * time.Hour)
}
