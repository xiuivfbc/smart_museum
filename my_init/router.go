package my_init

import (
	"github.com/gin-contrib/cors"
	"github.com/xiuivfbc/smart_museum/config"
	"github.com/xiuivfbc/smart_museum/controller"
	"github.com/xiuivfbc/smart_museum/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 初始化Gin路由
func SetupRouter() *gin.Engine {
	qrPath := config.Conf.GetString("server.path") + "/qrcodes"
	//设置路由
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	r.Static("/qrcodes", qrPath)

	// 静态文件服务 - 前端页面
	frontendPath := config.Conf.GetString("server.path") + "/frontend"
	r.Static("/frontend", frontendPath)

	// 根路径重定向到前端
	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/frontend/index.html")
	})

	// favicon 处理
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.Status(204) // No Content
	})

	r.GET("/ticket/use", controller.UseTicket)
	user := r.Group("/auth")
	{
		user.POST("/register", controller.RegisterUser)
		user.POST("/login", controller.LoginUser)
		user.POST("/upload", controller.UploadUser)
		user.POST("/refresh", controller.RefreshToken)
		user.POST("/logout", controller.Logout)
		user.POST("/verificationCode", controller.GetVerificationCode)
	}

	env := r.Group("/environments")
	{
		env.POST("", controller.GetAllEnvironmentsByRoom)
		env.POST("/name", controller.GetLastEnvironmentByRoom)
		env.DELETE("/name", controller.DeleteEnvironmentByRoom)
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
		device.POST("/get", controller.GetDeviceControl)
		device.POST("/update", controller.UpdateDeviceControl)
		device.POST("/autoenv", controller.AutoEnvironmentByDevice)
	}

	ai := r.Group("/ai")
	ai.Use(middleware.JWTAuthMiddleware())
	{
		ai.POST("/control", controller.AiControl)
	}
	return r
}
