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
	var req struct {
		Value string `json:"value"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	config.DeviceControl = req.Value

	// 向MQTT发布消息
	topic := config.Conf.GetString("mqtt.device_topic")
	if topic != "" {
		err := service.PublishMQTT(topic, req.Value)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "MQTT消息发送失败"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}
