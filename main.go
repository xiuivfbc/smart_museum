package main

import (
	"log"

	"github.com/xiuivfbc/smart_museum/my_init"
)

// main 作为程序入口，只负责启动服务和调用路由
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	port := my_init.GetServerPort()
	r := my_init.SetupRouter()

	log.Println("智能文化遗产保护后端服务启动，监听端口: ", port)
	r.Run(":" + port)
}
