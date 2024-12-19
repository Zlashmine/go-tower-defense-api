package cache

import (
	"time"

	"github.com/go-redis/redis/v8"
)

const CacheExpiryTime = time.Hour * 24

func NewRedisClient(addr, pw string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pw,
		DB:       db,
	})
}

func NewRedisStore(client *redis.Client) Store {
	return Store{
		Users: &UserStore{client},
	}
}
