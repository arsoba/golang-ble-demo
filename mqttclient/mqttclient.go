package mqttclient

import (
	"encoding/json"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttClient struct {
	topic          string
	PahoMqttClient mqtt.Client
}

func NewClient() MqttClient {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://172.20.10.5:1883")

	mqttClient := mqtt.NewClient(opts)
	return MqttClient{
		PahoMqttClient: mqttClient,
	}
}

func (m *MqttClient) Connect() {
	if token := m.PahoMqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func (m *MqttClient) Publish(topic string, payload interface{}) {
	message, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	token := m.PahoMqttClient.Publish(topic, byte(0), false, message)
	token.Wait()
}
