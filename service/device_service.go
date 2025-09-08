package service

import (
	"group_ten_server/config"
	"group_ten_server/dao"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetDeviceControlService(c *gin.Context, req interface{}) {
	reqStruct, ok := req.(struct {
		Device string `json:"device"`
	})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数类型错误"})
		return
	}
	log.Println("请求的设备:", reqStruct.Device)
	log.Println("当前设备控制状态:", config.DeviceControl[reqStruct.Device])
	if len(config.DeviceControl[reqStruct.Device]) != 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "设备状态异常"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"value": config.DeviceControl[reqStruct.Device]})
}

func UpdateDeviceControlService(c *gin.Context, req interface{}) {
	reqStruct, ok := req.(struct {
		Device string `json:"device"`
		Value  string `json:"value"`
	})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数类型错误"})
		return
	}
	log.Println("请求的设备:", reqStruct.Device, "更新值:", reqStruct.Value)
	if len(reqStruct.Value) != 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "value长度必须为10"})
		return
	}
	config.DeviceControl[reqStruct.Device] = reqStruct.Value
	log.Println("更新后设备控制状态:", config.DeviceControl[reqStruct.Device])

	topic := config.Conf.GetString("mqtt.device_topic")
	if topic != "" {
		err := PublishMQTT(topic, reqStruct.Value+reqStruct.Device[len(reqStruct.Device)-1:])
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "MQTT消息发送失败"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

func AutoEnvironmentByDeviceService(c *gin.Context, req interface{}) {
	reqStruct, ok := req.(struct {
		Device      string  `json:"device"`
		Temperature float64 `json:"温度"`
		Humidity    float64 `json:"湿度"`
	})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数类型错误"})
		return
	}
	room := "room" + reqStruct.Device[len(reqStruct.Device)-1:]
	log.Println("room:", room)
	nowenv, _ := dao.GetEnvironmentByNameFromTable(config.AppConfigInstance.RoomMapping[room])
	nowdev := []byte(config.DeviceControl[reqStruct.Device])
	log.Println("当前环境数据:", nowenv)

	if reqStruct.Temperature > nowenv.Temperature && config.DeviceControl[reqStruct.Device] != "0" {
		nowdev[0] = 0
		nowdev[1]++
	}
	if reqStruct.Temperature < nowenv.Temperature && config.DeviceControl[reqStruct.Device] != "3" {
		nowdev[0] = 0
		nowdev[1]--
	}
	if reqStruct.Humidity > nowenv.Humidity {
		nowdev[10] = 1
		nowdev[12] = byte(int(reqStruct.Humidity)/10 + '0')
		nowdev[13] = byte(int(reqStruct.Humidity)%10 + '0')
	}
	log.Println("更新后设备控制状态:", string(nowdev))
	config.DeviceControl[reqStruct.Device] = string(nowdev)
	topic := config.Conf.GetString("mqtt.device_topic")
	if topic != "" {
		err := PublishMQTT(topic, config.DeviceControl[reqStruct.Device]+reqStruct.Device[len(reqStruct.Device)-1:])
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "MQTT消息发送失败"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "自动环境处理完成"})
}
