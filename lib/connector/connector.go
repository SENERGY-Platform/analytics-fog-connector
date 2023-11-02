package connector

import (
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
