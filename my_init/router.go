package my_init

import (
	"group_ten_server/config"
	"group_ten_server/controller"
	"group_ten_server/middleware"

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
	r.GET("/ticket/use", controller.UseTicket)
	user := r.Group("/auth")
	{
		user.POST("/register", controller.RegisterUser)
		user.POST("/login", controller.LoginUser)
		user.POST("/upload", controller.UploadUser)
	}

	env := r.Group("/environments")
	{
		env.POST("", controller.GetAllEnvironmentsByRoom)
		env.POST(":name", controller.GetLastEnvironmentByRoom)
		env.DELETE(":name", controller.DeleteEnvironmentByRoom)
	}

	ticket := r.Group("/ticket")
	ticket.Use(middleware.JWTAuthMiddleware())
	{
		ticket.POST("/create", controller.CreateTicket)
		ticket.POST("/list", controller.ListTicket)
		ticket.DELETE("/delete", controller.DeleteTicket)
		ticket.PUT("/update", controller.UpdateTicket)

		ticket.GET("/entrycount/now", controller.GetTodayEntryCount)
		ticket.GET("/entrycount/all", controller.GetAllEntryCounts)
	}

	device := r.Group("/device")
	device.Use(middleware.JWTAuthMiddleware())
	{
		device.GET("/get", controller.GetDeviceControl)
		device.POST("/update", controller.UpdateDeviceControl)
	}
	return r
}
