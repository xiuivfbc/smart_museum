package controller

import (
	"group_ten_server/config"
	"group_ten_server/dao"
	"group_ten_server/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 获取DeviceControl值
func GetDeviceControl(c *gin.Context) {
	var req struct {
		Device string `json:"device"`
	}
	c.JSON(http.StatusOK, gin.H{"value": config.DeviceControl[req.Device]})
}

// 更新DeviceControl值
func UpdateDeviceControl(c *gin.Context) {
	var req struct {
		Device string `json:"device"`
		Value  string `json:"value"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	config.DeviceControl[req.Device] = req.Value

	// 向MQTT发布消息
	topic := config.Conf.GetString("mqtt.device_topic")
	if topic != "" {
		err := service.PublishMQTT(topic, req.Value+req.Device[len(req.Device)-1:])
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "MQTT消息发送失败"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

func AutoEnvironmentByDevice(c *gin.Context) {
	var req struct {
		Device string `json:"device"`
		Data   struct {
			Temperature float64 `json:"temperature"`
			Humidity    float64 `json:"humidity"`
		} `json:"data"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	room := "room" + req.Device[len(req.Device)-1:]
	nowenv, _ := dao.GetEnvironmentByNameFromTable(room)
	nowdev := []byte(config.DeviceControl[req.Device])

	if req.Data.Temperature > nowenv.Temperature && config.DeviceControl[req.Device] != "0" {
		nowdev[0] = 0
		nowdev[1]++
	}
	if req.Data.Temperature < nowenv.Temperature && config.DeviceControl[req.Device] != "3" {
		nowdev[0] = 0
		nowdev[1]--
	}
	if req.Data.Humidity > nowenv.Humidity {
		// 湿度高了，处理...
	}
	if req.Data.Humidity < nowenv.Humidity {
		// 湿度低了，处理...
	}
	config.DeviceControl[req.Device] = string(nowdev)
	// 向MQTT发布消息
	topic := config.Conf.GetString("mqtt.device_topic")
	if topic != "" {
		err := service.PublishMQTT(topic, config.DeviceControl[req.Device]+req.Device[len(req.Device)-1:])
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "MQTT消息发送失败"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "自动环境处理完成"})
}
