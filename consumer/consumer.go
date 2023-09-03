package consumer

import (
	"github.com/harshavardhan/event_delivery/config"
	"github.com/harshavardhan/event_delivery/redis"
	"log"
	"time"
)

func ConsumeEvents() {
	for {
		// can parallelize across destinations?
		for _, destination := range config.Destinations {
			log.Print("Processing events from " + destination)
			redis.ConsumeEvents(time.Now().UnixNano(), destination)
		}
		time.Sleep(5000 * time.Millisecond)
	}
}
