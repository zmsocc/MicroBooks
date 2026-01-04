package ioc

import (
	"context"
	"fmt"
	rlock "github.com/gotomicro/redis-lock"
	"github.com/redis/go-redis/v9"
	"time"
)

func InitRedis() redis.Cmdable {
	//addr := viper.GetString("redis.addr")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("Redis 连接失败：%v\n", err)
		return nil
	}
	return redisClient
}

func InitRLockClient(cmd redis.Cmdable) *rlock.Client {
	return rlock.NewClient(cmd)
}
