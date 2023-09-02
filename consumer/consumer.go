package consumer

import (
	"github.com/harshavardhan/event_delivery/config"
	"log"
	"time"
)

func ConsumeEvents() {
	for {
		// can parallelize across destinations?
		for _, destination := range config.Destinations {
			log.Print("Processing events from " + destination)
		}
		time.Sleep(5000 * time.Millisecond)
	}
}
