package storage

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	saveTime = 60 * time.Minute
)

type Rediscache struct {
	client *redis.Client
}

func NewRedisCache() *Rediscache {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	return &Rediscache{
		client: client,
	}
}

func (c Rediscache) get(expression string) []byte {
	result, _ := c.client.Get(context.Background(), expression).Bytes()
	return result
}

func (c Rediscache) Get() ([][]string, error) {
	ctx := context.Background()

	keys, err := c.client.LRange(ctx, "keylist", -10, -1).Result()
	if err != nil {
		return nil, err
	}

	history := make([][]string, 0, 9)
	for _, key := range keys {
		value := c.get(key)

		pair := []string{key, string(value)}
		history = append(history, pair)
	}

	return history, nil
}

func (c Rediscache) Set(expression, result string) error {
	ctx := context.Background()
	err := c.client.Set(ctx, expression, result, saveTime).Err()
	if err != nil {
		return err
	}

	err = c.client.RPush(ctx, "keylist", expression).Err()
	return err
}

func (c Rediscache) GetClient() *redis.Client {
	return c.client
}
