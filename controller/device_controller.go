package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xiuivfbc/smart_museum/service"
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
	service.GetDeviceControlService(c, req)
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
	service.UpdateDeviceControlService(c, req)
}

// AutoEnvironmentByDevice 根据设备数据自动调整环境
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
	service.AutoEnvironmentByDeviceService(c, req)
}
