package domain

import (
	"testing"
	"time"
)

func TestNewUser(t *testing.T) {
	// GIVEN
	firstName := "Vitalii"
	lastName := "Honchar"
	username := "vitalii_honchar"
	var chatId int64 = 12345
	now := time.Now()

	expectedNT := time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, now.Location())

	if expectedNT.Before(now) {
		expectedNT = now
	}

	// WHEN
	user := NewUser(firstName, lastName, username, chatId, now)

	// THEN
	if user.FirstName != firstName {
		t.Errorf("FirstName error: expected = %s, actual = %s\n", firstName, user.FirstName)
	}
	if user.LastName != lastName {
		t.Errorf("LastName error: expected = %s, actual = %s\n", lastName, user.LastName)
	}
	if user.Username != username {
		t.Errorf("Username error: expected = %s, actual = %s\n", username, user.Username)
	}
	if user.ChatId != chatId {
		t.Errorf("ChatId error: expected = %d, actual = %d\n", chatId, user.ChatId)
	}
	if user.NotificationTime != expectedNT {
		t.Errorf("NotificationTime error: expected = %v, actual = %v\n",
			expectedNT, user.NotificationTime)
	}
}

func TestUser_NextNotificationTime(t *testing.T) {
	// GIVEN
	user := NewUser("Vitalii", "Honchar", "vitalii_honchar", 12345, time.Now())
	nt := user.NotificationTime

	// WHEN
	user.NextNotificationTime()

	// THEN
	if user.NotificationTime.Sub(nt) != 24*time.Hour {
		t.Errorf("Next notification time should be after 24 hours: previous = %v, current = %v\n",
			nt, user.NotificationTime)
	}
}
