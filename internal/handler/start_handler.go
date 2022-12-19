package handler

import (
	"funny-bot/internal/database"
	"funny-bot/internal/domain"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"time"
)

type StartHandler struct {
	repository *database.UserRepository
}

func NewStartHandler(r *database.UserRepository) *StartHandler {
	return &StartHandler{repository: r}
}

const startMsg = `Hello, I'm Funny Bot. I will send you notifications about funny time.
So even you are working so hard, you will never forget about fun for yourself and your family!`

func (s *StartHandler) Matches(m *tgbotapi.Message) bool {
	return m.Text == "/start"
}

func (s *StartHandler) Handle(m *tgbotapi.Message) (*tgbotapi.MessageConfig, error) {
	u := newUser(m)
	if !s.repository.ExistsByChatId(u.ChatId) {
		s.repository.Save(*u)
	}

	log.Printf("Bot started: user = %+v\n", u)

	msg := tgbotapi.NewMessage(m.Chat.ID, startMsg)
	msg.ReplyToMessageID = m.MessageID

	return &msg, nil
}

func newUser(m *tgbotapi.Message) *domain.User {
	return &domain.User{
		Username:         m.Chat.UserName,
		FirstName:        m.Chat.FirstName,
		LastName:         m.Chat.LastName,
		ChatId:           m.Chat.ID,
		NotificationTime: time.Now(),
	}
}
