package scheduler

import (
	"funny-bot/internal/database"
	"funny-bot/internal/domain"
	"log"
	"sync"
	"time"
)

type NotificationScheduler struct {
	userRepository *database.UserRepository
	interval       time.Duration
	started        bool
	mutex          sync.Mutex
}

func NewNotificationScheduler(ur *database.UserRepository, d time.Duration) *NotificationScheduler {
	return &NotificationScheduler{
		userRepository: ur,
		interval:       d,
	}
}

func (ns *NotificationScheduler) Start() {
	ns.mutex.Lock()
	defer ns.mutex.Unlock()

	if !ns.started {
		ns.started = true
		go func() {
			timer := time.NewTimer(ns.interval)
			for t := range timer.C {
				log.Printf("Send notifications to funny users: time = %v\n", t)
				ns.SendNotifications()
				timer.Reset(ns.interval)
			}
		}()
	}
}

func (ns *NotificationScheduler) SendNotifications() {
	users := ns.userRepository.FindAllByNotificationTimeLessOrEquals(time.Now())

	var channels []<-chan bool

	for _, u := range users {
		channels = append(channels, ns.sendNotification(&u))
	}

	for _, c := range channels {
		<-c
	}
}

func (ns *NotificationScheduler) sendNotification(u *domain.User) <-chan bool {
	c := make(chan bool, 1)

	go func() {
		defer close(c)
		log.Printf("Send notification to user: %v\n", u)
		time.Sleep(5 * time.Second)
		log.Printf("Notification was sent to user: %v\n", u)
		c <- true
	}()
	return c
}
