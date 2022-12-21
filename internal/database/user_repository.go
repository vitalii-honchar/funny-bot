package database

import (
	"database/sql"
	"funny-bot/internal/domain"
	"time"
)

type UserRepository struct {
	users []domain.User
	db    *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) Save(u domain.User) {
	var updated bool
	for i, user := range ur.users {
		if user.ChatId == u.ChatId {
			ur.users[i] = u
			updated = true
			break
		}
	}
	if !updated {
		ur.users = append(ur.users, u)
	}
}

func (ur *UserRepository) FindByChatId(id int64) *domain.User {
	for _, u := range ur.users {
		if u.ChatId == id {
			return &u
		}
	}
	return nil
}

func (ur *UserRepository) ExistsByChatId(id int64) bool {
	return ur.FindByChatId(id) != nil
}

func (ur *UserRepository) FindAllByNotificationTimeLessOrEquals(t time.Time) []domain.User {
	var res []domain.User
	for _, u := range ur.users {
		if u.NotificationTime.Before(t) || u.NotificationTime == t {
			res = append(res, u)
		}
	}
	return res
}
