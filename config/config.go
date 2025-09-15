package config

import (
	"log"

	"github.com/spf13/viper"
)

type AppConfig struct {
	RoomMapping map[string]string `mapstructure:"room_mapping"`
}

var Conf *viper.Viper
var AppConfigInstance AppConfig
var Device1 string = "0000000001"
var Device2 string = "0000000001"
var Device3 string = "0000000001"
var Device4 string = "0000000001"
var DeviceControl = make(map[string]string)

const (
	TOTAL_NUM = "total_num"
)

// InitConfig 初始化配置文件读取
func InitConfig() {
	Conf = viper.New()
	Conf.SetConfigName("config")
	Conf.SetConfigType("yaml")
	Conf.AddConfigPath("./config")
	err := Conf.ReadInConfig()
	if err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}
	// 将配置解析到结构体
	err = Conf.Unmarshal(&AppConfigInstance)
	if err != nil {
		log.Fatalf("解析配置到结构体失败: %v", err)
	}

	DeviceControl["device1"] = Device1
	DeviceControl["device2"] = Device2
	DeviceControl["device3"] = Device3
	DeviceControl["device4"] = Device4

	log.Printf("配置文件加载成功，共加载 %d 个房间映射", len(AppConfigInstance.RoomMapping))
}
