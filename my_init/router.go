package my_init

import (
	"group_ten_server/controller"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

// SetupRouter 初始化Gin路由
func SetupRouter() *gin.Engine {
	//设置路由
	r := gin.Default()
	r.Use(cors.Default())
	user := r.Group("/auth")
	{
		user.POST("/register", controller.Register)
		user.POST("/login", controller.Login)
	}

	// 环境数据相关路由
	env := r.Group("/environments")
	{
		env.POST("", controller.CreateEnvironment)
		env.GET("", controller.GetAllEnvironments)
		env.GET(":name", controller.GetEnvironmentByName)
		env.PUT(":name", controller.UpdateEnvironmentByName)
		env.DELETE(":name", controller.DeleteEnvironmentByName)
	}

	ticket := r.Group("/ticket")
	{
		ticket.POST("/create", controller.CreateTicket)
		ticket.GET("/list", controller.ListTickets)
		ticket.DELETE("/delete", controller.DeleteTicket)
		ticket.PUT("/update", controller.UpdateTicket)
	}
	return r
}
