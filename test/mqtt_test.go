package test

import (
	"group_ten_server/config"
	"group_ten_server/service"
	"testing"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func TestMQTTConnectAndPubSub(t *testing.T) {
	config.InitConfig()
	service.InitMQTT()

	topic := "test/topic"
	msg := "hello mqtt"
	ch := make(chan string, 1)

	// 订阅
	err := service.SubscribeMQTT(topic, func(client MQTT.Client, message MQTT.Message) {
		ch <- string(message.Payload())
	})
	if err != nil {
		t.Fatalf("订阅失败: %v", err)
	}

	time.Sleep(1 * time.Second) // 等待订阅生效

	// 发布
	err = service.PublishMQTT(topic, msg)
	if err != nil {
		t.Fatalf("发布失败: %v", err)
	}

	select {
	case recv := <-ch:
		if recv != msg {
			t.Errorf("接收到的消息不一致，期望: %s, 实际: %s", msg, recv)
		}
	case <-time.After(3 * time.Second):
		t.Error("超时未收到消息")
	}
}
