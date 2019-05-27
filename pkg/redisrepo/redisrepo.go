package redisrepo

import "github.com/go-redis/redis"

type RedisRepository struct {
	client *redis.Client
}

func NewRepository(database int, password string) *RedisRepository {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // use default Addr
		Password: password,         // no password set
		DB:       database,         // use default DB
	})
	return &RedisRepository{client}
}
