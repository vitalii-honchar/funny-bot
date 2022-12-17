package handler

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type MessageHandler interface {
	Handle(m *tgbotapi.Message) (*tgbotapi.MessageConfig, error)

	Matches(m *tgbotapi.Message) bool
}
