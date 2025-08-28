package my_init

import "group_ten_server/config"

func GetServerPort() (port string) {
	if s := config.Conf.GetString("gin.port"); s != "" {
		port = s
	} else {
		port = "8080"
	}
	return port
}
