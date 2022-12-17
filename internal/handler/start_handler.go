package handler

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type StartHandler struct{}

func (s *StartHandler) Matches(m *tgbotapi.Message) bool {
	return m.Text == "/start"
}

func (s *StartHandler) Handle(m *tgbotapi.Message) (*tgbotapi.MessageConfig, error) {
	msg := tgbotapi.NewMessage(m.Chat.ID, "Start triggered")
	msg.ReplyToMessageID = m.MessageID

	return &msg, nil
}
