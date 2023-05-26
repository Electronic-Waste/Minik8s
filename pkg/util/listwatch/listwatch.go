package listwatch

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type WatchHandler func(msg *redis.Message)

type WatchResult struct {
	ObjectType string //Pod, Deployment
	ActionType string //apply, delete, update
	Payload    []byte //struct pod or deployment
}

// TODO(shaowang): Expand to multiple machines in the future
var rdb = redis.NewClient(&redis.Options{
	Addr:     "192.168.1.6:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

var sub *redis.PubSub = nil

func Subscribe(topic string) <-chan *redis.Message {
	print("redis: subscribe " + topic + "\n")
	sub = rdb.Subscribe(ctx, topic)
	return sub.Channel()
}

func Unsubscribe(topic string) error {
	if sub == nil {
		return fmt.Errorf(
			"Should first subscribe a topic",
		)
	}
	err := sub.Unsubscribe(ctx, topic)
	if err != nil {
		return fmt.Errorf(
			"Unsubscribe from channel %s falied!",
			topic,
		)
	}
	return nil
}

func Publish(topic string, msg interface{}) {
	print("redis: publish " + topic + "\n")
	rdb.Publish(ctx, topic, msg)
}

// When using this function, you should add "go" keyword in front of it.
func Watch(topic string, handler WatchHandler) {
	channel := Subscribe(topic)
	for msg := range channel {
		fmt.Println("redis: receive msg")
		handler(msg)
	}
}
