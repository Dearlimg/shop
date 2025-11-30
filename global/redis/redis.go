package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client
var ctx = context.Background()

// InitRedis 初始化Redis连接
func InitRedis(addr, password string, db int) error {
	Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// 测试连接
	_, err := Client.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	log.Println("Redis connection established")
	return nil
}

// CloseRedis 关闭Redis连接
func CloseRedis() error {
	if Client != nil {
		return Client.Close()
	}
	return nil
}

// GetClient 获取Redis客户端
func GetClient() *redis.Client {
	return Client
}

// GetContext 获取上下文
func GetContext() context.Context {
	return ctx
}
