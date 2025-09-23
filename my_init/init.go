package my_init

import (
	"context"

	"github.com/xiuivfbc/smart_museum/config"
	"github.com/xiuivfbc/smart_museum/dao"
	"github.com/xiuivfbc/smart_museum/middleware"
	"github.com/xiuivfbc/smart_museum/service"
)

func init() {
	config.InitConfig()

	middleware.InitJwtKey(config.Conf.GetString("gin.jwt_secret"))
	service.InitJwtKey(config.Conf.GetString("gin.jwt_secret"))

	service.InitMQTT()
	// CollectDataFromHardware(config.Conf.GetString("mqtt.environment_topic"))

	dao.InitDB()
	dao.EnsureActivationCodes()

	InitRedis()
	service.RedisClient.Set(context.Background(), config.TOTAL_NUM, 0, 0)
	StartDailyEntryCountJob()

	// GetMaxKBPower()
}
