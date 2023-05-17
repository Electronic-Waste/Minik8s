package listwatch

import (
	"context"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type WatchHandler func(msg *redis.Message)

// TODO(shaowang): Expand to multiple machines in the future
var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

func Subscribe(topic string) <-chan *redis.Message {
	print("redis: subscribe " + topic + "\n")
	sub := rdb.Subscribe(ctx, topic)
	return sub.Channel()
}

func Publish(topic string, msg interface{}) {
	print("redis: publish " + topic + "\n")
	rdb.Publish(ctx, topic, msg)
}

func Watch(topic string, handler WatchHandler) {
	print("redis: watch " + topic + "\n")
	channel := Subscribe(topic)
	for true{
		for msg := range channel {
			print("redis: handle msg\n")
			handler(msg)
		}
	}
	
}
