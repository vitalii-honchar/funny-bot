package main

import (
	"funny-bot/internal/app"
	"funny-bot/internal/database"
	"funny-bot/internal/handler"
	"funny-bot/internal/scheduler"
	"funny-bot/internal/telegram"
	"log"
	"os"
	"time"
)

type config struct {
	botToken string
}

func readConfig() config {
	return config{
		botToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
	}
}

func main() {
	cfg := readConfig()
	bot, err := telegram.NewBot(cfg.botToken)
	if err != nil {
		log.Fatalln(err)
	}
	repository := database.NewUserRepository()
	fs := app.NewFunnyService(repository, bot)
	ns := scheduler.NewNotificationScheduler(fs, time.Second)

	bot.AddHandler(handler.NewStartHandler(repository))
	ns.Start()
	bot.Start()
}
