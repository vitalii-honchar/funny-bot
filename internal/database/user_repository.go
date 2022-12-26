package database

import (
	"database/sql"
	"funny-bot/internal/domain"
	"funny-bot/internal/lib"
	"funny-bot/internal/time_provider"
	"log"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

const selectByChatId = "SELECT * FROM tg_user WHERE chat_id = $1"
const existsByChatId = "SELECT id FROM tg_user WHERE chat_id = $1"
const insertUser = "INSERT INTO tg_user (first_name, last_name, username, chat_id, notification_time) VALUES ($1, $2, $3, $4, $5)" +
	"ON CONFLICT (chat_id)" +
	"DO UPDATE SET first_name = $1, last_name = $2, username = $3, notification_time = $5"
const selectByNotificationTime = "SELECT * FROM tg_user WHERE notification_time <= $1 ORDER BY notification_time"

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) Save(u *domain.User) <-chan bool {
	return lib.Async(func(c chan bool) {
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
	return lib.Async(func(c chan *domain.User) {
		r := ur.db.QueryRow(selectByChatId, id)
		u, err := scanUser(r)
		if err != nil {
			c <- nil
		} else {
			c <- u
		}
	})
}

func (ur *UserRepository) ExistsByChatId(id int64) <-chan bool {
	return lib.Async(func(c chan bool) {
		r := ur.db.QueryRow(existsByChatId, id)
		var id int64
		if err := r.Scan(&id); err != nil {
			c <- false
		} else {
			c <- true
		}
	})
}

func (ur *UserRepository) FindAllByNotificationTimeLessOrEquals(t time.Time) <-chan []*domain.User {
	return lib.Async(func(c chan []*domain.User) {
		rows, err := ur.db.Query(selectByNotificationTime, t.In(time.UTC))
		if err != nil {
			log.Printf("Unexpected error during read list of users: notificationTime = %v, err = %v\n", t, err)
			c <- make([]*domain.User, 0)
		} else {
			defer rows.Close()
			var users []*domain.User
			for rows.Next() {
				u, _ := scanUser(rows)
				if u != nil {
					users = append(users, u)
				}
			}
			c <- users
		}
	})
}

func scanUser(r userScanner) (*domain.User, error) {
	var result domain.User
	var nt time.Time
	err := r.Scan(&result.Id, &result.FirstName, &result.LastName,
		&result.Username, &result.ChatId, &nt)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("Can't scan result: err = %v\n", err)
		}
		return nil, err
	}
	result.NotificationTime = nt.In(time_provider.EstLocation)
	return &result, nil
}

type userScanner interface {
	Scan(...any) error
}
