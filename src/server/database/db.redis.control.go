package database

type RedisStore struct{}

func (c *RedisStore) Set(key string, value string) error {
	return Instance.redisClient.Set(key, value, 0).Err()
}

func (c *RedisStore) Get(key string) (string, error) {
	return Instance.redisClient.Get(key).Result()
}

func (c *RedisStore) Del(key string) error {
	return Instance.redisClient.Get(key).Err()
}
