package redis

import (
	"context"
	"fmt"

	lib "github.com/redis/go-redis/v9"
)

type Redis struct {
	client *lib.Client
	//pubsub *lib.
}

func NewRedis(connUrl string) *Redis {
	opts, err := lib.ParseURL(connUrl)
	if err != nil {
		panic(err)
	}

	return &Redis{client: lib.NewClient(opts)}
}

func (r *Redis) Publish(ctx context.Context, topic, message string) error {
	msgID, err := r.client.XAdd(ctx, &lib.XAddArgs{
		Stream: "quorum",
		// Use * for auto-generating an ID
		ID: "*",
		Values: map[string]any{
			"field": "value",
			"foo":   "bar",
		},
	}).Result()

	if err != nil {
		panic(err)
	}
	fmt.Println("Added message with ID:", msgID)
	return nil
}
