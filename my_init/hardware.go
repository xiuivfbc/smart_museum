package my_init

import (
	"encoding/json"
	"group_ten_server/config"
	"group_ten_server/controller"
	"group_ten_server/model"
	"group_ten_server/service"
	"log"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// 从硬件收集数据
func CollectDataFromHardware(topic string) {
	go func() {
		// 只需订阅一次，回调会持续处理
		err := service.SubscribeMQTT(topic, func(client MQTT.Client, msg MQTT.Message) {
			log.Println("收到MQTT消息:", string(msg.Payload()))
			// 解析完整json字段
			var data struct {
				Room        string  `json:"room"`
				Temperature float64 `json:"temperature"`
				Humidity    float64 `json:"humidity"`
				Gas         int     `json:"gas"`
				AutoFan     int     `json:"auto_fan"`
				FanLevel    int     `json:"fan_level"`
				AutoLight   int     `json:"auto_light"`
				Light       int     `json:"light"`
				Human       int     `json:"human"`
				Dark        int     `json:"dark"`
			}
			if err := json.Unmarshal(msg.Payload(), &data); err != nil {
				log.Println("解析环境数据失败:", err)
				return
			}
			// 保存环境数据到对应房间表
			if data.Room == "" {
				log.Println("缺少 room 字段或格式错误")
				data.Room = "room1" // 默认房间
			}
			nowdev := []byte(config.DeviceControl["device"+data.Room[len(data.Room)-1:]])
			log.Println("当前设备控制状态:", string(nowdev))
			nowdev[0] = byte(data.AutoFan) + 48
			nowdev[1] = byte(data.FanLevel) + 48
			nowdev[3] = byte(data.AutoLight) + 48
			nowdev[4] = byte(data.Light) + 48
			config.DeviceControl["device"+data.Room[len(data.Room)-1:]] = string(nowdev)
			log.Println("更新后设备控制状态:", config.DeviceControl["device"+data.Room[len(data.Room)-1:]])

			env := model.Environment{
				Temperature: data.Temperature,
				Humidity:    data.Humidity,
				Gas:         data.Gas,
				AutoFan:     data.AutoFan,
				FanLevel:    data.FanLevel,
				AutoLight:   data.AutoLight,
				Light:       data.Light,
				Human:       data.Human,
				Dark:        data.Dark,
			}
			controller.CreateEnvironmentByRoom(data.Room, &env)
		})
		if err != nil {
			log.Println("订阅主题失败:", err)
		}
	}()
}
