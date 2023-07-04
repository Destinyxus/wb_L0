package store

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	redis *redis.Client
}

func NewCache() *Cache {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "10120001",
		DB:       0,
	})
	return &Cache{
		redis: client,
	}
}

func (c *Cache) InitCache() error {
	pong, err := c.redis.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println(pong)
	}

	return nil
}

func (c *Cache) SaveToCache(order []byte, id string) error {

	err := c.redis.Set(context.Background(), "orderID:"+id, order, 0).Err()
	if err != nil {
		log.Fatal(err)
	}
	return nil

}

func (c *Cache) GetOrderByID(id string) ([]byte, bool) {

	order, err := c.redis.Get(context.Background(), "orderID:"+id).Bytes()
	if err == redis.Nil {
		return nil, false

	} else if err != nil {
		log.Fatal(err)
	}
	return order, true
}
