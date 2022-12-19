package main

import (
	"funny-bot/internal/app"
	"funny-bot/internal/database"
	"funny-bot/internal/handler"
	"funny-bot/internal/scheduler"
	"funny-bot/internal/telegram"
	"log"
	"time"
)

const funnyMessage = "Funny time starts! Stop your deals and go to have a fun!"

func main() {
	bot, err := telegram.NewBot("5633412468:AAFiQ4H_CJOflrFJaimGZj6LdR3NmO8xWyw")
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
