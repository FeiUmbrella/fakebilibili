package redis

import (
	"context"
	"fakebilibili/infrastructure/config"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/sethvargo/go-retry"
	"log"
	"time"
)

var RedisDb *redis.Client

func init() {
	redisConf := config.Config.RedisConfig
	b := retry.NewFibonacci(10 * time.Second) // 设置斐波那契基数
	ctx := context.Background()
	// 连接失败后，按照斐波那契设置退避时间，重连5次
	if err := retry.Do(ctx, retry.WithMaxRetries(5, b), func(ctx context.Context) error {
		RedisDb = redis.NewClient(&redis.Options{
			Addr:         fmt.Sprintf("%s:%d", redisConf.IP, redisConf.Port),
			Password:     redisConf.Password,
			DB:           0, //共有8个db，use default DB 0号
			PoolSize:     4, // 最大连接数，默认为4*cpu个数 todo:上服务器后要调大
			MinIdleConns: 5, //最少活跃连接数
		})
		_, err := RedisDb.Ping().Result() // 看是否连通
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		// 多次重连仍失败
		log.Fatalf("重试5次后仍然无法连接redis，请排查redis服务端是否启动/配置信息是否正确，错误信息为： %v \n", err)
	}
}
