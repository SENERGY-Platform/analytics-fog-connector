package connector

import (
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/mqtt"
)

type Connector struct {
	FogMQTTClient *mqtt.MQTTClient
	PlatformMQTTClient *mqtt.MQTTClient
	PublishResultsToPlatform bool
	UserID string
}

func NewConnector(fogMqttClient *mqtt.MQTTClient, platformMqttClient *mqtt.MQTTClient, publishResultsToPlatform bool, userID string) *Connector {
	return &Connector{
		FogMQTTClient: fogMqttClient,
		PlatformMQTTClient: platformMqttClient,
		PublishResultsToPlatform: publishResultsToPlatform,
		UserID: userID,
	}
}
