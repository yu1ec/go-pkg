package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yu1ec/go-pkg/cachex/driver"
)

func init() {
	driver.Register("redis", New)
}

type RedisDriver struct {
	ctx    context.Context
	client *redis.Client
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func New(config any) (driver.Driver, error) {
	cfg, ok := config.(*RedisConfig)
	if !ok {
		return nil, fmt.Errorf("invalid config for redis cache")
	}

	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// 检查连接是否成功
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisDriver{ctx: ctx, client: client}, nil
}

// 实现 BaseDriver 接口
func (r *RedisDriver) Add(k string, v any, d time.Duration) error {
	success, err := r.client.SetNX(r.ctx, k, v, d).Result()
	if err != nil {
		return err
	}
	if !success {
		return fmt.Errorf("key already exists")
	}
	return nil
}

func (r *RedisDriver) Delete(k string) {
	r.client.Del(r.ctx, k)
}

func (r *RedisDriver) DeleteExpired() {
	// Redis 自动处理过期键，无需手动实现
}

func (r *RedisDriver) Flush() {
	r.client.FlushAll(r.ctx)
}

func (r *RedisDriver) Get(k string) (any, bool) {
	val, err := r.client.Get(r.ctx, k).Result()
	if err == redis.Nil {
		return nil, false
	}
	return val, err == nil
}

func (r *RedisDriver) GetWithExpiration(k string) (any, time.Time, bool) {
	val, err := r.client.Get(r.ctx, k).Result()
	if err == redis.Nil {
		return nil, time.Time{}, false
	}
	if err != nil {
		return nil, time.Time{}, false
	}
	ttl, err := r.client.TTL(r.ctx, k).Result()
	if err != nil {
		return val, time.Time{}, true
	}
	if ttl < 0 {
		return val, time.Time{}, true
	}
	return val, time.Now().Add(ttl), true
}

func (r *RedisDriver) Replace(k string, x any, d time.Duration) error {
	_, err := r.client.Get(r.ctx, k).Result()
	if err == redis.Nil {
		return fmt.Errorf("key not found")
	}
	return r.client.Set(r.ctx, k, x, d).Err()
}

func (r *RedisDriver) Set(k string, x any, d time.Duration) {
	r.client.Set(r.ctx, k, x, d)
}

func (r *RedisDriver) SetDefault(k string, x any) {
	r.client.Set(r.ctx, k, x, 0)
}

// 实现 NumericOperations 接口
func (r *RedisDriver) IncrementInt(k string, n int) (int, error) {
	return int(r.client.IncrBy(r.ctx, k, int64(n)).Val()), nil
}

func (r *RedisDriver) DecrementInt(k string, n int) (int, error) {
	return int(r.client.DecrBy(r.ctx, k, int64(n)).Val()), nil
}

func (r *RedisDriver) IncrementInt64(k string, n int64) (int64, error) {
	return r.client.IncrBy(r.ctx, k, n).Result()
}

func (r *RedisDriver) DecrementInt64(k string, n int64) (int64, error) {
	return r.client.DecrBy(r.ctx, k, n).Result()
}

func (r *RedisDriver) IncrementUint(k string, n uint) (uint, error) {
	val, err := r.client.IncrBy(r.ctx, k, int64(n)).Uint64()
	return uint(val), err
}

func (r *RedisDriver) DecrementUint(k string, n uint) (uint, error) {
	val, err := r.client.DecrBy(r.ctx, k, int64(n)).Uint64()
	return uint(val), err
}

func (r *RedisDriver) IncrementUint64(k string, n uint64) (uint64, error) {
	return r.client.IncrBy(r.ctx, k, int64(n)).Uint64()
}

func (r *RedisDriver) DecrementUint64(k string, n uint64) (uint64, error) {
	return r.client.DecrBy(r.ctx, k, int64(n)).Uint64()
}
