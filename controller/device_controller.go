package controller

import (
	"group_ten_server/config"
	"group_ten_server/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 获取DeviceControl值
func GetDeviceControl(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"value": config.DeviceControl})
}

// 更新DeviceControl值
func UpdateDeviceControl(c *gin.Context) {
	req := c.Query("value")
	config.DeviceControl = req

	// 向MQTT发布消息
	topic := config.Conf.GetString("mqtt.topic")
	if topic != "" {
		err := service.PublishMQTT(topic, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "MQTT消息发送失败"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}
