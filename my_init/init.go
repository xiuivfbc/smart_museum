package my_init

import (
	"group_ten_server/config"
	"group_ten_server/controller"
	"group_ten_server/dao"
	"group_ten_server/middleware"
	"group_ten_server/service"
)

func init() {
	config.InitConfig()
	middleware.InitJwtKey(config.Conf.GetString("gin.jwt_secret"))
	controller.InitJwtKey(config.Conf.GetString("gin.jwt_secret"))
	service.InitMQTT()
	dao.InitDB()
	dao.EnsureActivationCodes()
	CollectDataFromHardware(config.Conf.GetString("mqtt.environment_topic"))
	StartDailyEntryCountJob()
	GetMaxKBPower()
}
