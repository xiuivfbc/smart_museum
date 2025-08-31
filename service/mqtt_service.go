package service

import (
	"encoding/json"
	"group_ten_server/config"
	"group_ten_server/controller"
	"group_ten_server/model"
	"log"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var mqttClient MQTT.Client

// InitMQTT 初始化MQTT客户端（支持连接任意MQTT Broker，包括EMQX等）
func InitMQTT() {
	broker := config.Conf.GetString("mqtt.broker")
	clientID := config.Conf.GetString("mqtt.client_id")
	if broker == "" {
		broker = "tcp://127.0.0.1:1883"
	}
	if clientID == "" {
		clientID = "go_smart_home_client"
	}
	opts := MQTT.NewClientOptions().AddBroker(broker)
	opts.SetUsername("cp")
	opts.SetPassword(config.Conf.GetString("mqtt.client_password"))
	opts.SetClientID(clientID)
	mqttClient = MQTT.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

// 从硬件收集数据
func CollectDataFromHardware(topic string) {
	go func() {
		// 只需订阅一次，回调会持续处理
		err := SubscribeMQTT(topic, func(client MQTT.Client, msg MQTT.Message) {
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

// PublishMQTT 发布消息到指定主题
func PublishMQTT(topic, payload string) error {
	token := mqttClient.Publish(topic, 0, false, payload)
	token.Wait()
	return token.Error()
}

// SubscribeMQTT 订阅主题并处理消息
func SubscribeMQTT(topic string, callback MQTT.MessageHandler) error {
	token := mqttClient.Subscribe(topic, 0, callback)
	token.Wait()
	return token.Error()
}
