package router

import (
	"github.com/gin-contrib/cors"
	"group_ten_server/controller"
	"group_ten_server/dao"

	"github.com/gin-gonic/gin"
)

// SetupRouter 初始化Gin路由
func SetupRouter() *gin.Engine {
	// 数据库连接配置
	// 初始化数据库连接
	dao.InitDB()
	// 可通过 dao.GetDB() 获取 *sql.DB 实例

	//设置路由
	r := gin.Default()
	r.Use(cors.Default())
	user := r.Group("/auth")
	{
		user.POST("/register", controller.Register)
		user.POST("/login", controller.Login)
	}
	return r
}
