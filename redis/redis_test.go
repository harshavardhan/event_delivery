package redis

import (
	"github.com/harshavardhan/event_delivery/models"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestStoreEvent(t *testing.T) {
	Init()

	ev := models.Event{
		UserID:  "1",
		Payload: "data1",
	}
	now := time.Now().UnixNano()
	id := StoreEvent("A", now, ev)

	var em models.EventMetadata
	_ = mapstructure.WeakDecode(redisClient.HGetAll(ctx, id).Val(), &em)

	assert.Equal(t, em.UserID, ev.UserID, "userId in datastore is different")
	assert.Equal(t, em.Payload, ev.Payload, "payload in datastore is different")

	Cleanup()
}

func TestProcessEvent(t *testing.T) {
	Init()

	inIds := make([]string, 0)
	outIds := make([]string, 0)

	n := 100
	for i := 1; i <= n; i++ {
		ev := models.Event{
			UserID:  strconv.Itoa(rand.Intn(i)),
			Payload: strconv.Itoa(i),
		}
		id := StoreEvent("A", time.Now().UnixNano(), ev)
		inIds = append(inIds, id)
	}

	for len(outIds) != n {
		processedIds := ProcessEvents(time.Now().UnixNano(), "A")
		outIds = append(outIds, processedIds...)
		time.Sleep(100 * time.Millisecond)
	}
	assert.Equal(t, inIds, outIds, "Events are processed in a different order")

	Cleanup()
}
