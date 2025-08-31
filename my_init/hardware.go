package my_init

import (
	"encoding/json"
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
			// 先解析出 room 字段
			var rawData map[string]interface{}
			json.Unmarshal(msg.Payload(), &rawData)
			room, ok := rawData["room"].(string)
			if !ok {
				log.Println("缺少 room 字段或格式错误")
				return
			}
			delete(rawData, "room")
			jsonData, _ := json.Marshal(rawData)
			var env model.Environment
			if err := json.Unmarshal(jsonData, &env); err != nil {
				log.Println("解析环境数据失败:", err)
				return
			}
			controller.CreateEnvironmentByRoom(room, &env)
		})
		if err != nil {
			log.Println("订阅主题失败:", err)
		}
	}()
}
