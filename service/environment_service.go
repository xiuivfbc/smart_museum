package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xiuivfbc/smart_museum/config"
	"github.com/xiuivfbc/smart_museum/dao"
)

func GetLastEnvironmentByRoomService(c *gin.Context, req interface{}) {
	reqStruct, ok := req.(struct {
		Room string `json:"room"`
	})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数类型错误"})
		return
	}
	table, ok := config.AppConfigInstance.RoomMapping[reqStruct.Room]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效房间名"})
		return
	}
	env, err := dao.GetEnvironmentByNameFromTable(table)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "未找到"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": env})
}

func GetAllEnvironmentsByRoomService(c *gin.Context, req interface{}) {
	reqStruct, ok := req.(struct {
		Room string `json:"room"`
	})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数类型错误"})
		return
	}
	table, ok := config.AppConfigInstance.RoomMapping[reqStruct.Room]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效房间名"})
		return
	}
	envs, err := dao.GetAllEnvironmentsFromTable(table)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": envs})
}

func DeleteEnvironmentByRoomService(c *gin.Context, req interface{}) {
	reqStruct, ok := req.(struct {
		Room string `json:"room"`
	})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数类型错误"})
		return
	}
	table, ok := config.AppConfigInstance.RoomMapping[reqStruct.Room]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效房间名"})
		return
	}
	if err := dao.DeleteEnvironmentByNameFromTable(table); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
