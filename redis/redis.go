package redis

import (
	"context"
	"github.com/google/uuid"
	"github.com/harshavardhan/event_delivery/config"
	"github.com/harshavardhan/event_delivery/models"
	"github.com/harshavardhan/event_delivery/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"
	"time"
)

var redisClient *redis.Client
var ctx = context.Background()

// sorted queue key
func sqKey(destination string) string {
	return "sortedSet" + "$" + destination
}

func StoreEvent(ev models.Event) {
	var em = models.EventMetadata{
		Timestamp: time.Now().UnixNano(),
		UserID:    ev.UserID,
		Payload:   ev.Payload,
	}
	em.ExecTimestamp = em.Timestamp

	// Add broadcast to destinations
	for _, destination := range config.Destinations {
		id := uuid.New().String()
		// log.Printf("Adding %s to destination %s", id, destination)

		// need to handle redis errors later
		// Store event data mapped to id in a hash
		redisClient.HSet(ctx, id, em)

		// Each destination has a sorted set from which events are picked up by earliest time first
		redisClient.ZAdd(ctx, sqKey(destination), redis.Z{
			Score:  float64(em.ExecTimestamp),
			Member: utils.BuildKey(em.Timestamp, id),
		})

		// Each destination has a list for order in which events have to be processed
		redisClient.LPush(ctx, destination, id)
	}
}

func ConsumeEvents(before int64, destination string) {
	// need to add some element count limits here while fetching
	ids := redisClient.ZRangeByScore(ctx, sqKey(destination), &redis.ZRangeBy{
		Min: "0",
		Max: strconv.FormatInt(before, 10),
	}).Val()
	for _, cid := range ids {
		id := utils.GetId(cid)
		firstId := redisClient.LIndex(ctx, destination, -1).Val()

		execute := id == firstId
		successResponse := utils.MockSuccess()
		// log.Println(id, execute, successResponse)

		metadataMap := redisClient.HGetAll(ctx, id).Val()
		var em models.EventMetadata
		_ = mapstructure.WeakDecode(metadataMap, &em)

		if execute && successResponse {
			// might need to use multi and exec together here to update in a transaction
			redisClient.Del(ctx, id)
			redisClient.ZRem(ctx, sqKey(destination), utils.BuildKey(em.Timestamp, id))
			redisClient.RPop(ctx, destination)

			log.Printf("Successfully sent payload %s for destination %s", em.Payload, destination)
			continue
		}

		if !execute {
			// not the first-in one, so set execTime to the execTime of the first-in one
			firstExecTimestamp := utils.StrToInt(redisClient.HGet(ctx, firstId, "execTimestamp").Val())
			em.ExecTimestamp = firstExecTimestamp
		} else {
			// no success response from destination case
			// exponential backoff depending on retryCount
			em.ExecTimestamp = time.Now().UnixNano() + (1<<em.RetryCount)*int64(config.Delta)
			em.RetryCount += 1
		}
		// update metadata
		redisClient.HSet(ctx, id, em)
		// update exec time score
		redisClient.ZAdd(ctx, sqKey(destination), redis.Z{
			Score:  float64(em.ExecTimestamp),
			Member: utils.BuildKey(em.Timestamp, id),
		})
	}
}

func Init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: "",
		DB:       0,
	})
}
