package app

import (
	"funny-bot/internal/database"
	"funny-bot/internal/domain"
	"funny-bot/internal/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type FunnyService struct {
	repository *database.UserRepository
	bot        *telegram.Bot
}

const funnyMessage = "It's time for funny things! Stop everything and go to fun (write code in Kotlin :))"

func NewFunnyService(r *database.UserRepository, b *telegram.Bot) *FunnyService {
	return &FunnyService{
		repository: r,
		bot:        b,
	}
}

func (fs *FunnyService) SendNotifications() <-chan bool {
	res := make(chan bool, 1)

	go func() {
		defer close(res)
		users := fs.repository.FindAllByNotificationTimeLessOrEquals(CurrentTime())

		var channels []<-chan bool

		for _, u := range users {
			channels = append(channels, fs.sendNotification(&u))
		}

		for _, c := range channels {
			<-c
		}

		res <- true
	}()
	return res
}

func (fs *FunnyService) sendNotification(u *domain.User) <-chan bool {
	c := make(chan bool, 1)

	go func() {
		defer close(c)
		log.Printf("Send notification to user: %v\n", u)
		msg := tgbotapi.NewMessage(u.ChatId, funnyMessage)
		fs.bot.Send(&msg)
		u.NextNotificationTime()
		fs.repository.Save(*u)
		log.Printf("Notification was sent to user: %v\n", u)
		c <- true
	}()
	return c
}
