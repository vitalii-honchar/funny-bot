package bot

import (
	"funny-bot/internal/handler"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type Bot struct {
	bot      *tgbotapi.BotAPI
	handlers []handler.MessageHandler
}

func NewBot(token string) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &Bot{bot: bot}, nil
}

func (b *Bot) AddHandler(h handler.MessageHandler) *Bot {
	b.handlers = append(b.handlers, h)
	return b
}

func (b *Bot) Start() {
	log.Println("Starting Telegram bot...")
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)
	b.handleUpdates(updates)
}

func (b *Bot) handleUpdates(u tgbotapi.UpdatesChannel) {
	for update := range u {
		if update.Message != nil { // If we got a message
			b.handleUpdate(&update)
		}
	}
}

func (b *Bot) handleUpdate(u *tgbotapi.Update) {
	log.Printf("[%s] %s", u.Message.From.UserName, u.Message.Text)

	for _, h := range b.handlers {
		if h.Matches(u.Message) {
			resp, err := h.Handle(u.Message)
			if err == nil {
				_, err = b.bot.Send(resp)
			}
			if err != nil {
				log.Printf("Error during processing update: update = %v, error = %v\n", u, err)
			}
		}
	}
}
