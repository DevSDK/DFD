package database

import (
	"time"
)

//RedisStore data structure for redis access
type RedisStore struct {
	BaseDB
}

//SetWithExpire set value with expire duration
func (c *RedisStore) SetWithExpire(key string, value interface{}, expire time.Duration) error {
	return Instance.redisClient.Set(key, value, expire).Err()
}

//Set value without expire
func (c *RedisStore) Set(key string, value interface{}) error {
	return Instance.redisClient.Set(key, value, 0).Err()
}

//Get value by key
func (c *RedisStore) Get(key string) (string, error) {
	return Instance.redisClient.Get(key).Result()
}

//Del delete value by key
func (c *RedisStore) Del(key string) error {
	return Instance.redisClient.Get(key).Err()
}
