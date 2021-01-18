package database

import (
	"time"
)

type RedisStore struct {
	BaseDB
}

func (c *RedisStore) SetWithExpire(key string, value interface{}, expire time.Duration) error {
	return Instance.redisClient.Set(key, value, expire).Err()
}

func (c *RedisStore) Set(key string, value interface{}) error {
	return Instance.redisClient.Set(key, value, 0).Err()
}

func (c *RedisStore) Get(key string) (string, error) {
	return Instance.redisClient.Get(key).Result()
}

func (c *RedisStore) Del(key string) error {
	return Instance.redisClient.Get(key).Err()
}
