package time_provider

import (
	"log"
	"time"
)

var EstLocation *time.Location

func init() {
	l, err := time.LoadLocation("EST")
	if err != nil {
		log.Fatalln(err)
	}
	EstLocation = l
}

func CurrentTime() time.Time {
	return time.Now().In(EstLocation)
}
