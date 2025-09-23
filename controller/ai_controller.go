package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/xiuivfbc/smart_museum/service"
)

func AiControl(c *gin.Context) {
	var req struct {
		Question string `json:"question"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, gin.H{"error": "参数错误"})
		return
	}
	service.AiControlService(c, req)
}
