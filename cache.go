package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v8"
)

// 业务前缀
var Prefix = "demo"

// RedisClient redis客户端
type RedisClient struct {
	ctx    context.Context
	cache  *redis.Client
	locker *redislock.Client
}

// NewRedisClient 实例化redis服务
func NewRedisClient(ctx context.Context) *RedisClient {
	sentinel := NewSentinelnstance(ctx, &RedisSentinelConfig{
		MasterName:    "Sentinel4php",
		SentinelAddrs: strings.Split("127.0.0.1:26379,127.0.0.2:26379,127.0.0.3:26379", ","),
		Password:      "redis123456",
		DB:            10,
		ReadTimeout:   3,
		WriteTimeout:  1,
		IdleTimeout:   30,
	})
	cluster := NewClusterlnstance(ctx, &RedisClusterlConfig{
		Addrs:    strings.Split("127.0.0.1:26379,127.0.0.2:26379,127.0.0.3:26379", ","),
		Password: "redis123456",
	})
	return &RedisClient{
		ctx:    ctx,
		cache:  If(sentinel == nil, cluster, sentinel).(*redis.Client),
		locker: redislock.New(sentinel),
	}
}

// GetRedis 获取redis实例
func (c *RedisClient) GetRedis() *redis.Client {
	return c.cache
}

// GetLock 获取分布式锁
func (c *RedisClient) GetLock(key string, t time.Duration) (*redislock.Lock, error) {
	lock, err := c.locker.Obtain(c.ctx, key, t, nil)
	if err != nil {
		return nil, err
	}
	return lock, nil
}

// UnLock 释放锁
func (c *RedisClient) UnLock(lock *redislock.Lock) error {
	if lock == nil {
		return nil
	}
	return lock.Release(c.ctx)
}

// SaveStruct 缓存结构体
func (c *RedisClient) SaveStruct(k string, d interface{}, second time.Duration) error {
	jsonStr, err := json.Marshal(d)
	if err != nil {
		return fmt.Errorf("save key `%s` json marshal error", k)
	}
	err = c.cache.Set(c.ctx, k, string(jsonStr), time.Second*second).Err()
	if err != nil {
		return fmt.Errorf("save key `%s` fail:%s", k, err.Error())
	}
	return nil
}

// GetStruct 获取结构体
func (c *RedisClient) GetStruct(k string, data interface{}) error {
	valStr, err := c.cache.Get(c.ctx, k).Result()
	if redis.Nil == err {
		return fmt.Errorf("key `%s` not found", k)
	}
	if err != nil {
		return fmt.Errorf("get key `%s` fail:%s", k, err.Error())
	}
	return json.Unmarshal([]byte(valStr), &data)
}
