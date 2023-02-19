package dao

import (
	"context"
	"github.com/a754962942/project-project/config"
	"github.com/go-redis/redis/v8"
	"time"
)

var Rc *RedisCache

type RedisCache struct {
	rdb *redis.Client
}

func init() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.C.R.Addr,
		Password: config.C.R.Password,
		DB:       config.C.R.DB,
	})
	Rc = &RedisCache{
		rdb: rdb,
	}
}
func (r *RedisCache) Put(ctx context.Context, key, value string, expire time.Duration) error {
	err := r.rdb.Set(ctx, key, value, expire).Err()
	return err
}
func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	result, err := r.rdb.Get(ctx, key).Result()
	return result, err
}
