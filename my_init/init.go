package my_init

import (
	"context"
	"group_ten_server/config"
	"group_ten_server/dao"
	"group_ten_server/middleware"
	"group_ten_server/service"
)

func init() {
	config.InitConfig()

	middleware.InitJwtKey(config.Conf.GetString("gin.jwt_secret"))
	service.InitJwtKey(config.Conf.GetString("gin.jwt_secret"))

	service.InitMQTT()
	CollectDataFromHardware(config.Conf.GetString("mqtt.environment_topic"))

	dao.InitDB()
	dao.EnsureActivationCodes()

	InitRedis()
	service.RedisClient.Set(context.Background(), config.TOTAL_NUM, 0, 0)
	StartDailyEntryCountJob()

	GetMaxKBPower()
}
