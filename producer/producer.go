package producer

import (
	"bytes"
	"encoding/json"
	"github.com/harshavardhan/event_delivery/models"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func ProduceEvents() {
	client := &http.Client{}
	for i := 1; i <= 100; i++ {
		ev := models.Event{
			UserID:  strconv.Itoa(rand.Intn(i)),
			Payload: strconv.Itoa(i),
		}

		for {
			var buf bytes.Buffer
			_ = json.NewEncoder(&buf).Encode(ev)
			req, _ := http.NewRequest(http.MethodPost, "http://localhost:8090/receive_event", &buf)

			resp, _ := client.Do(req)

			if resp.StatusCode == http.StatusOK {
				break
			}
			log.Printf("Sleeping for producer %d", i)
			time.Sleep(1 * time.Second)
		}
	}
}
