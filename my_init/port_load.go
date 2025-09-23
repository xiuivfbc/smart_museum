package my_init

import "github.com/xiuivfbc/smart_museum/config"

func GetServerPort() (port string) {
	if s := config.Conf.GetString("gin.port"); s != "" {
		port = s
	} else {
		port = "8080"
	}
	return port
}
