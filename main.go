package main

import (
	"encoding/json"
	"github.com/harshavardhan/event_delivery/consumer"
	"github.com/harshavardhan/event_delivery/models"
	"github.com/harshavardhan/event_delivery/redis"
	"log"
	"net/http"
)

func parseRequest(req *http.Request) (ev models.Event) {
	// No validation on request method or body yet
	decoder := json.NewDecoder(req.Body)
	_ = decoder.Decode(&ev)
	log.Printf("%+v\n", ev)
	return
}

func sendResponse(w http.ResponseWriter, msg string) {
	var resp = models.EventResponse{
		Msg: msg,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func receiveEvent(w http.ResponseWriter, req *http.Request) {
	ev := parseRequest(req)
	redis.StoreEvent(ev)
	sendResponse(w, "Event received")
}

func serverInit() {
	http.HandleFunc("/receive_event", receiveEvent)

	port := "8090"
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Unable to start http service : %s", err)
		return
	}
}

func main() {
	redis.RedisInit()

	go func() { serverInit() }()

	consumer.ConsumeEvents()
}
