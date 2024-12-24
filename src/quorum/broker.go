package quorum

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// must be created using NewRedis()
type Broker struct {
	client *redis.Client
}

func NewRedis(connUrl string) *Broker {
	opts, err := redis.ParseURL(connUrl)
	if err != nil {
		panic(err)
	}

	return &Broker{client: redis.NewClient(opts)}
}

func (r *Broker) Write(ctx context.Context, topic string, id string, msg map[string]any) error {
	_, err := r.client.XAdd(ctx, &redis.XAddArgs{
		Stream: topic,
		ID:     id,
		Values: msg,
	}).Result()
	if err != nil {
		return fmt.Errorf("failed to publish, %s", err)
	}
	return nil
}

func (r *Broker) Read(ctx context.Context, topic string, batch int, timeout time.Duration) ([]map[string]any, error) {
	streams, err := r.client.XRead(ctx, &redis.XReadArgs{
		Streams: []string{topic, "$"},
		Count:   int64(batch),
		Block:   timeout,
	}).Result()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("failed to publish, %s", err)
	}
	if len(streams) != 1 {
		return nil, fmt.Errorf("unexpected stream count, %d", len(streams))
	}
	res := []map[string]any{}
	for _, message := range streams[0].Messages {
		resp := message.Values
		resp["msgId"] = message.ID
		res = append(res, resp)
	}
	return res, nil
}
