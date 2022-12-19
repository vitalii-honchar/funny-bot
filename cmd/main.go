package main

import (
	bot2 "funny-bot/internal/bot"
	"funny-bot/internal/database"
	"funny-bot/internal/handler"
	"funny-bot/internal/scheduler"
	"log"
	"time"
)

const funnyMessage = "Funny time starts! Stop your deals and go to have a fun!"

func main() {
	bot, err := bot2.NewBot("5633412468:AAFiQ4H_CJOflrFJaimGZj6LdR3NmO8xWyw")
	if err != nil {
		log.Fatalln(err)
	}
	repository := database.NewUserRepository()
	ns := scheduler.NewNotificationScheduler(repository, time.Second)

	bot.AddHandler(handler.NewStartHandler(repository))
	ns.Start()
	bot.Start()
}
