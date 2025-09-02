package my_init

import (
	"context"
	"group_ten_server/config"
	"group_ten_server/controller"

	"github.com/redis/go-redis/v9"
)

func InitRedis() {
	controller.RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.Conf.GetString("redis.addr"),
		Password: config.Conf.GetString("redis.password"),
		DB:       config.Conf.GetInt("redis.db"),
	})
	// 可选：测试连接
	ctx := context.Background()
	if err := controller.RedisClient.Ping(ctx).Err(); err != nil {
		panic("Redis连接失败: " + err.Error())
	}
}
