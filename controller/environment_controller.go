package controller

import (
	"fmt"
	"group_ten_server/config"
	"group_ten_server/dao"
	"group_ten_server/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetLastEnvironmentByRoom(c *gin.Context) {
	room := c.Query("room")
	table, ok := config.AppConfigInstance.RoomMapping[room]
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

func GetAllEnvironmentsByRoom(c *gin.Context) {
	room := c.Query("room")
	table, ok := config.AppConfigInstance.RoomMapping[room]
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

// CreateEnvironmentByRoom 支持内部调用，参数为env和room
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
	room := c.Query("room")
	table, ok := config.AppConfigInstance.RoomMapping[room]
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
