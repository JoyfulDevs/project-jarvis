package app

import (
	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	rdb *redis.Client
}

func NewRedisService(addr string, pw string) *RedisService {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pw,
		DB:       0,
	})
	return &RedisService{
		rdb: client,
	}
}

func (r *RedisService) Close() error {
	return r.rdb.Close()
}
