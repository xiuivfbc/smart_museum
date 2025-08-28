package main

import (
	"group_ten_server/router"
	"log"
)

// main 作为程序入口，只负责启动服务和调用路由
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	r := router.SetupRouter()
	log.Println("智能文化遗产保护后端服务启动，监听端口: 8080")
	r.Run(":8080")
}
