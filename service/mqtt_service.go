package service

import (
	"group_ten_server/config"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var mqttClient MQTT.Client

// InitMQTT 初始化MQTT客户端（支持连接任意MQTT Broker，包括EMQX等）
func InitMQTT() {
	broker := config.Conf.GetString("mqtt.broker")
	clientID := config.Conf.GetString("mqtt.client_id")

	opts := MQTT.NewClientOptions().AddBroker(broker)
	opts.SetUsername("cp")
	opts.SetPassword(config.Conf.GetString("mqtt.client_password"))
	opts.SetClientID(clientID)
	mqttClient = MQTT.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
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
