package time_provider

import (
	"log"
	"time"
)

var location *time.Location

func init() {
	l, err := time.LoadLocation("EST")
	if err != nil {
		log.Fatalln(err)
	}
	location = l
}

func CurrentTime() time.Time {
	return time.Now().In(location)
}
