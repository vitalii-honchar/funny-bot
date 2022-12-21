package app

import (
	"funny-bot/internal/database"
	"funny-bot/internal/domain"
	"funny-bot/internal/lib"
	"funny-bot/internal/telegram"
	"funny-bot/internal/time_provider"
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
	return lib.Async(func(res chan bool) {
		users := fs.repository.FindAllByNotificationTimeLessOrEquals(time_provider.CurrentTime())

		var channels []<-chan bool

		for _, u := range users {
			channels = append(channels, fs.sendNotification(&u))
		}

		for _, c := range channels {
			<-c
		}

		res <- true
	})
}

func (fs *FunnyService) sendNotification(u *domain.User) <-chan bool {
	return lib.Async(func(c chan bool) {
		log.Printf("Send notification to user: %v\n", u)
		msg := tgbotapi.NewMessage(u.ChatId, funnyMessage)
		fs.bot.Send(&msg)
		u.NextNotificationTime()
		fs.repository.Save(u)
		log.Printf("Notification was sent to user: %v\n", u)
		c <- true
	})
}
