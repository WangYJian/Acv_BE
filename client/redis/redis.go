package redis

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client
var ctx = context.Background()

// 初始化 Redis 客户端
func InitRedis() {
	// 从环境变量中读取 Redis 相关配置
	redisClient = redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     os.Getenv("REDIS_ADDRESS") + ":" + os.Getenv("REDIS_PORT"), // Redis 服务器地址和端口号
		Password: os.Getenv("REDIS_PASSWORD"),                                // Redis 认证密码
		DB:       0,                                                          // Redis 数据库编号
	})
	fmt.Println(redisClient)
}

func Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return redisClient.Set(ctx, key, value, expiration)
}
func Get(key string) *redis.StringCmd {
	return redisClient.Get(ctx, key)
}
func Del(keys string) *redis.IntCmd {
	return redisClient.Del(ctx, keys)
}
func Exists(keys string) *redis.IntCmd {
	return redisClient.Exists(ctx, keys)
}
func HIncrBy(key string, field string, incr int64) *redis.IntCmd {
	return redisClient.HIncrBy(ctx, key, field, incr)
}
