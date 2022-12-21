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
	botToken        string
	databaseConnUrl string
}

func readConfig() config {
	return config{
		botToken:        os.Getenv("TELEGRAM_BOT_TOKEN"),
		databaseConnUrl: os.Getenv("DB_CONNECTION_URL"),
	}
}

func main() {
	cfg := readConfig()
	bot, err := telegram.NewBot(cfg.botToken)
	if err != nil {
		log.Fatalln(err)
	}
	db, err := database.OpenConnection(cfg.databaseConnUrl)
	if err != nil {
		log.Fatalln(err)
	}
	repository := database.NewUserRepository(db)
	fs := app.NewFunnyService(repository, bot)
	ns := scheduler.NewNotificationScheduler(fs, time.Second)

	bot.AddHandler(handler.NewStartHandler(repository))
	ns.Start()
	bot.Start()
}
