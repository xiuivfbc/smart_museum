package config

import (
	"log"

	"github.com/spf13/viper"
)

var Conf *viper.Viper

// InitConfig 初始化配置文件读取
func InitConfig() {
	Conf = viper.New()
	Conf.SetConfigName("config")
	Conf.SetConfigType("yaml")
	Conf.AddConfigPath("../config")
	err := Conf.ReadInConfig()
	if err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}
}
