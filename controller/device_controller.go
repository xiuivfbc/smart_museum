package controller

import (
	"group_ten_server/config"
	"group_ten_server/dao"
	"group_ten_server/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 获取DeviceControl值
func GetDeviceControl(c *gin.Context) {
	var req struct {
		Device string `json:"device"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	log.Println("请求的设备:", req.Device)
	log.Println("当前设备控制状态:", config.DeviceControl[req.Device])
	if len(config.DeviceControl[req.Device]) != 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "设备状态异常"})
		return
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
	log.Println("请求的设备:", req.Device, "更新值:", req.Value)
	if len(req.Value) != 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "value长度必须为10"})
		return
	}
	config.DeviceControl[req.Device] = req.Value
	log.Println("更新后设备控制状态:", config.DeviceControl[req.Device])

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
			Temperature float64 `json:"温度"`
			Humidity    float64 `json:"湿度"`
		} `json:"data"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	room := "room" + req.Device[len(req.Device)-1:]
	log.Println("room:", room)
	nowenv, _ := dao.GetEnvironmentByNameFromTable(room)
	nowdev := []byte(config.DeviceControl[req.Device])
	log.Println("当前环境数据:", nowenv)

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
	log.Println("更新后设备控制状态:", string(nowdev))
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
