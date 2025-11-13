package core

import (
	"blog/internal/global"
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var ctx = context.Background()

func InitRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     global.Config.Redis.Addr,
		Password: global.Config.Redis.Password,
		DB:       global.Config.Redis.DB,
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		logrus.Fatalf("redis连接失败,%v", err)
		return nil
	}
	logrus.Infof("redis连接成功")
	return rdb
}
