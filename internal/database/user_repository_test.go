package database

import (
	"funny-bot/internal/domain"
	"math/rand"
	"testing"
	"time"
)

var random *rand.Rand

func init() {
	s1 := rand.NewSource(time.Now().UnixNano())
	random = rand.New(s1)
}

func TestUserRepository_Save(t *testing.T) {
	repository := newTestRepository(t)
	user := domain.NewUser("first_name", "last_name", "user_name", random.Int63(), time.Now())

	saved := <-repository.Save(user)

	if !saved {
		t.Errorf("Save result is unexpected: saved = %v\n", saved)
	}

	su := <-repository.FindByChatId(user.ChatId)
	if su == nil {
		t.Fatalf("Can't find user by chatId: chatId = %d\n", user.ChatId)
	}
	if su.Id == 0 {
		t.Errorf("Saved id should be greater than 0: id = %d\n", su.Id)
	}
	if su.FirstName != user.FirstName {
		t.Errorf("Saved first name is wrong: expected = %s, actual = %s\n", user.FirstName, su.FirstName)
	}
	if su.LastName != user.LastName {
		t.Errorf("Saved last name is wrong: expected = %s, actual = %s\n", user.LastName, su.LastName)
	}
	if su.Username != user.Username {
		t.Errorf("Saved username is wrong: expected = %s, actual = %s\n", user.Username, su.Username)
	}
	if su.ChatId != user.ChatId {
		t.Errorf("Saved chat id is wrong: expected = %d, actual = %d\n", user.ChatId, su.ChatId)
	}
	if !su.NotificationTime.Equal(user.NotificationTime) {
		t.Errorf("Saved notification time wrong: expected = %v, actual = %v\n",
			user.NotificationTime, su.NotificationTime)
	}
}

func newTestRepository(t *testing.T) *UserRepository {
	db := openTestConnection(t)
	return NewUserRepository(db)
}
