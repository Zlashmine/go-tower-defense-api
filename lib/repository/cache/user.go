package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"

	"tower-defense-api/lib/models"
)

type UserStore struct {
	client *redis.Client
}

func (store *UserStore) Get(ctx context.Context, id int64) (*models.User, error) {
	cacheKey := fmt.Sprintf("user:%d", id)

	if store.client.Exists(ctx, cacheKey).Val() == 0 {
		return nil, nil
	}

	val, err := store.client.Get(ctx, cacheKey).Result()

	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var user *models.User

	if val != "" {
		if err := json.Unmarshal([]byte(val), &user); err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (store *UserStore) Set(ctx context.Context, user *models.User) error {
	json, err := json.Marshal(user)

	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("user:%d", user.ID)

	return store.client.SetEX(
		ctx,
		cacheKey,
		json,
		CacheExpiryTime,
	).
		Err()
}
