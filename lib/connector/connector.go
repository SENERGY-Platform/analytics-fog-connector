package connector

import (
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/mqtt"
)

type Connector struct {
	FogMQTTClient *mqtt.MQTTClient
	PlatformMQTTClient *mqtt.MQTTClient
}

func NewConnector(fogMqttClient *mqtt.MQTTClient, platformMqttClient *mqtt.MQTTClient) *Connector {
	return &Connector{
		FogMQTTClient: fogMqttClient,
		PlatformMQTTClient: platformMqttClient,
	}
}
