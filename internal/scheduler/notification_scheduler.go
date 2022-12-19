package scheduler

import (
	"funny-bot/internal/app"
	"sync"
	"time"
)

type NotificationScheduler struct {
	funnyService *app.FunnyService
	interval     time.Duration
	started      bool
	mutex        sync.Mutex
}

func NewNotificationScheduler(fs *app.FunnyService, d time.Duration) *NotificationScheduler {
	return &NotificationScheduler{
		funnyService: fs,
		interval:     d,
	}
}

func (ns *NotificationScheduler) Start() {
	ns.mutex.Lock()
	defer ns.mutex.Unlock()

	if !ns.started {
		ns.started = true
		go func() {
			timer := time.NewTimer(ns.interval)
			for range timer.C {
				<-ns.funnyService.SendNotifications()
				timer.Reset(ns.interval)
			}
		}()
	}
}
