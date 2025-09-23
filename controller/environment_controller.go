package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xiuivfbc/smart_museum/config"
	"github.com/xiuivfbc/smart_museum/dao"
	"github.com/xiuivfbc/smart_museum/model"
	"github.com/xiuivfbc/smart_museum/service"
)

func GetLastEnvironmentByRoom(c *gin.Context) {
	var req struct {
		Room string `json:"room"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	service.GetLastEnvironmentByRoomService(c, req)
}

func GetAllEnvironmentsByRoom(c *gin.Context) {
	var req struct {
		Room string `json:"room"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	service.GetAllEnvironmentsByRoomService(c, req)
}

// CreateEnvironmentByRoom 仅支持内部调用，参数为env和room
func CreateEnvironmentByRoom(room string, env *model.Environment) error {
	table, ok := config.AppConfigInstance.RoomMapping[room]
	if !ok {
		return fmt.Errorf("无效房间名")
	}
	if err := dao.CreateEnvironmentToTable(table, env); err != nil {
		return err
	}
	return nil
}

func DeleteEnvironmentByRoom(c *gin.Context) {
	var req struct {
		Room string `json:"room"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	service.DeleteEnvironmentByRoomService(c, req)
}
