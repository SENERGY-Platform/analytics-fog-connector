package connector

import (
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/control"
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/mqtt"
)

type Connector struct {
	MQTTClient *mqtt.MQTTClient
}

func NewConnector(mqttClient *mqtt.MQTTClient) *Connector {
	return &Connector{
		MQTTClient: mqttClient,
	}
}

func (controller *Connector) ForwardToMaster(message []byte) {
	controller.MQTTClient.Publish(control.ControlTopic, string(message), 2)
}
