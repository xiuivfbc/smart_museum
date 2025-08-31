package my_init

import (
	"group_ten_server/config"
	"group_ten_server/controller"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

// SetupRouter 初始化Gin路由
func SetupRouter() *gin.Engine {
	qrPath := config.Conf.GetString("server.path") + "/qrcodes"
	//设置路由
	r := gin.Default()
	r.Use(cors.Default())
	r.Static("/qrcodes", qrPath)
	user := r.Group("/auth")
	{
		user.POST("/register", controller.RegisterUser)
		user.POST("/login", controller.LoginUser)
		user.POST("/upload", controller.UploadUser)
	}

	// 环境数据相关路由
	env := r.Group("/environments")
	{
		env.GET("", controller.GetAllEnvironmentsByRoom)
		env.GET(":name", controller.GetLastEnvironmentByRoom)
		env.DELETE(":name", controller.DeleteEnvironmentByRoom)
	}

	ticket := r.Group("/ticket")
	{
		ticket.POST("/create", controller.CreateTicket)
		ticket.GET("/list", controller.ListTicket)
		ticket.DELETE("/delete", controller.DeleteTicket)
		ticket.PUT("/update", controller.UpdateTicket)
		ticket.GET("/use", controller.UseTicket)
	}
	return r
}
