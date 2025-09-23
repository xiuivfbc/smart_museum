package my_init

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/xiuivfbc/smart_museum/config"
	"github.com/xiuivfbc/smart_museum/service"
)

func InitRedis() {
	service.RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.Conf.GetString("redis.addr"),
		Password: config.Conf.GetString("redis.password"),
		DB:       config.Conf.GetInt("redis.db"),
	})
	// 可选：测试连接
	ctx := context.Background()
	if err := service.RedisClient.Ping(ctx).Err(); err != nil {
		panic("Redis连接失败: " + err.Error())
	}
}
