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
