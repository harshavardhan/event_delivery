package main

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"time"
)

var redisClient *redis.Client

// constant list of possible destinations, move to constants/config later
var destinations = []string{"A", "B", "C", "D"}

type event struct {
	UserID  string `json:"userId"`
	Payload string
}

type eventResponse struct {
	Msg string `json:"msg"`
}

type eventMetadata struct {
	// duplicating event fields to avoid embedded structs reformat before write to redis
	UserID        string `redis:"userId"`
	Payload       string `redis:"payload"`
	Timestamp     int64  `redis:"timestamp"`     // (to be generally set by the client) set to server time if not sent (unix epoch nanoseconds)
	ExecTimestamp int64  `redis:"execTimestamp"` // initially equal to timestamp, on retries or waits set to next exec timestamp
	RetryCount    int    `redis:"retryCount"`
}

func parseRequest(req *http.Request) (ev event) {
	// No validation on request method or body yet
	decoder := json.NewDecoder(req.Body)
	_ = decoder.Decode(&ev)
	log.Printf("%+v\n", ev)
	return
}

func sendResponse(w http.ResponseWriter, msg string) {
	var resp = eventResponse{
		Msg: msg,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func storeEvent(ev event) {
	var em = eventMetadata{
		Timestamp: time.Now().UnixNano(),
		UserID:    ev.UserID,
		Payload:   ev.Payload,
	}
	em.ExecTimestamp = em.Timestamp

	// Add event metadata to redis
	id := uuid.New().String()
	log.Println(id)
	ctx := context.Background()
	// need to handle redis errors later
	// Store event data mapped to id in a hash
	redisClient.HSet(ctx, id, em)

	// Add broadcast to destinations
	for _, destination := range destinations {
		// Each destination has a sorted set from which events are picked up by earliest time first
		redisClient.ZAdd(ctx, destination, redis.Z{
			Score:  float64(em.ExecTimestamp),
			Member: id,
		})

		// Each destination has a list for order in which events have to be processed
		redisClient.LPush(ctx, destination, id)
	}
}

func receiveEvent(w http.ResponseWriter, req *http.Request) {
	ev := parseRequest(req)
	storeEvent(ev)
	sendResponse(w, "Event received")
}

func main() {
	http.HandleFunc("/receive_event", receiveEvent)

	redisAddr := "localhost:6379"
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	port := "8090"
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Unable to start http service : %s", err)
		return
	}
}
