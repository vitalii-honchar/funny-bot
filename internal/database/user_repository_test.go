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
	user := newTestUser()

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

func TestUserRepository_ExistsByChatId(t *testing.T) {
	repository := newTestRepository(t)

	t.Run("exists should be true", func(t *testing.T) {
		user := newTestUser()

		<-repository.Save(user)

		exists := <-repository.ExistsByChatId(user.ChatId)
		if !exists {
			t.Error("Unexpected exists result: exists should be true")
		}
	})

	t.Run("exists should be false", func(t *testing.T) {
		user := newTestUser()

		exists := <-repository.ExistsByChatId(user.ChatId)
		if exists {
			t.Error("Unexpected exists result: exists should be false")
		}
	})
}

func newTestRepository(t *testing.T) *UserRepository {
	db := openTestConnection(t)
	return NewUserRepository(db)
}

func newTestUser() *domain.User {
	return domain.NewUser("first_name", "last_name", "user_name", random.Int63(), time.Now())
}
