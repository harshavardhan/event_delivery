package redis

import (
	"github.com/harshavardhan/event_delivery/models"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
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
