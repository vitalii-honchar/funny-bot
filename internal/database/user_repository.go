package database

import (
	"database/sql"
	"funny-bot/internal/domain"
	"funny-bot/internal/time_provider"
	"log"
	"time"
)

type UserRepository struct {
	users []domain.User
	db    *sql.DB
}

const selectByChatId = "SELECT * FROM tg_user WHERE chat_id = $1"
const insertUser = "INSERT INTO tg_user (first_name, last_name, username, chat_id, notification_time) VALUES ($1, $2, $3, $4, $5)"

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) Save(u *domain.User) <-chan bool {
	return async(func(c chan bool) {
		_, err := ur.db.Exec(insertUser, u.FirstName, u.LastName, u.Username, u.ChatId, u.NotificationTime.In(time.UTC))
		if err != nil {
			log.Printf("Unexpected error during saving of user: user = %v, err = %v\n", u, err)
			c <- false
		} else {
			c <- true
		}
	})
}

func (ur *UserRepository) FindByChatId(id int64) <-chan *domain.User {
	return async(func(c chan *domain.User) {
		var result domain.User
		r := ur.db.QueryRow(selectByChatId, id)
		var nt time.Time
		err := r.Scan(&result.Id, &result.FirstName, &result.LastName,
			&result.Username, &result.ChatId, &nt)
		if err != nil {
			if err != sql.ErrNoRows {
				log.Printf("Can't scan result: chatId = %d, err = %v\n", id, err)
			}
			c <- nil
		} else {
			result.NotificationTime = nt.In(time_provider.EstLocation)
			c <- &result
		}
	})
}

func (ur *UserRepository) ExistsByChatId(id int64) <-chan bool {
	return async(func(c chan bool) {
		u := <-ur.FindByChatId(id)
		c <- u != nil
	})
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

func async[V any](f func(chan V)) <-chan V {
	c := make(chan V, 1)
	go func() {
		defer close(c)
		f(c)
	}()
	return c
}
