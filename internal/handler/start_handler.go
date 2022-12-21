package handler

import (
	"funny-bot/internal/database"
	"funny-bot/internal/domain"
	"funny-bot/internal/time_provider"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type StartHandler struct {
	repository *database.UserRepository
}

func NewStartHandler(r *database.UserRepository) *StartHandler {
	return &StartHandler{repository: r}
}

const startMsg = `Hello, I'm Funny Bot. 
I will send you notifications about funny time every day at 6 p.m. EST (I'm still in development mode, I will be more flexible).
So even you are working so hard, you will never forget about fun for yourself and your family!`

func (s *StartHandler) Matches(m *tgbotapi.Message) bool {
	return m.Text == "/start"
}

func (s *StartHandler) Handle(m *tgbotapi.Message) (*tgbotapi.MessageConfig, error) {
	u := newUser(m)
	exists := <-s.repository.ExistsByChatId(u.ChatId)
	if !exists {
		<-s.repository.Save(u)
	}

	log.Printf("Bot started: user = %+v\n", u)

	msg := tgbotapi.NewMessage(m.Chat.ID, startMsg)
	msg.ReplyToMessageID = m.MessageID

	return &msg, nil
}

func newUser(m *tgbotapi.Message) *domain.User {
	return domain.NewUser(
		m.Chat.FirstName,
		m.Chat.LastName,
		m.Chat.UserName,
		m.Chat.ID,
		time_provider.CurrentTime())
}
