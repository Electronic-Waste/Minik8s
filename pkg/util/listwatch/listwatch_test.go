package listwatch

import (
	"testing"
	"github.com/go-redis/redis/v8"
)

func TestListwatch(t *testing.T) {
	go Watch(
		"test",
		func (msg *redis.Message) {
			t.Logf("watcher receive %s", msg.Payload)
		},
	)
	Publish("test", "1111")
	Unsubscribe("test")
	Publish("test", "2222")
}