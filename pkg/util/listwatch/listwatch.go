package listwatch

import (
    "context"
    "github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type WatchHandler func(msg *redis.Message)

// TODO(shaowang): Expand to multiple machines in the future
var rdb = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
	Password: "", // no password set
	DB:	0,	// use default DB
})

func Subscribe(topic string) (<-chan *redis.Message){
	sub := rdb.Subscribe(ctx, topic)
	return sub.Channel()
}

func Publish(topic string, msg interface{}) {
	rdb.Publish(ctx, topic, msg)
}

func Watch(topic string, handler WatchHandler) {
	channel := Subscribe(topic)
	for msg := range channel {
		handler(msg)
	}
}