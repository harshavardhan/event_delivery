package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

var redisClient *redis.Client

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

func redisInit() {
	redisAddr := "localhost:6379"
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})
}
