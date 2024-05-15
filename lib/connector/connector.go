package connector

import (
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/clients/mqtt"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/itf"
)

type Connector struct {
	FogMQTTClient mqtt.MQTTClient
	PlatformMQTTClient mqtt.MQTTClient
	PublishResultsToPlatform bool
	UserID string
	LocalMessageRelayHandler itf.MessageRelayHandler
	CloudMessageRelayHandler itf.MessageRelayHandler

}

func NewConnector(fogMqttClient mqtt.MQTTClient, platformMqttClient mqtt.MQTTClient, publishResultsToPlatform bool, userID string, LocalMessageRelayHandler itf.MessageRelayHandler, CloudMessageRelayHandler itf.MessageRelayHandler) *Connector {
	return &Connector{
		FogMQTTClient: fogMqttClient,
		PlatformMQTTClient: platformMqttClient,
		PublishResultsToPlatform: publishResultsToPlatform,
		UserID: userID,
		LocalMessageRelayHandler: LocalMessageRelayHandler,
		CloudMessageRelayHandler: CloudMessageRelayHandler,
	}
}
