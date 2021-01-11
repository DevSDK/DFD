package database

func (c *DBInstance) SetRedis(key string, value string) error {
	return c.redisClient.Set(key, value, 0).Err()
}

func (c *DBInstance) GetRedis(key string) (string, error) {
	return c.redisClient.Get(key).Result()
}

func (c *DBInstance) DelRedis(key string) error {
	return c.redisClient.Get(key).Err()
}
