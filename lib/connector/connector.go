package connector

import (
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/clients/mqtt"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/itf"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/logging"
	"sync"
)

type Connector struct {
	FogMQTTClient mqtt.MQTTClient
	PlatformMQTTClient mqtt.MQTTClient
	PublishResultsToPlatform bool
	UserID string
	LocalMessageRelayHandler itf.MessageRelayHandler
	CloudMessageRelayHandler itf.MessageRelayHandler
	SubscriptedLocalTopics map[string]struct{}
	mu                      sync.RWMutex
}

func NewConnector(fogMqttClient mqtt.MQTTClient, platformMqttClient mqtt.MQTTClient, publishResultsToPlatform bool, userID string, LocalMessageRelayHandler itf.MessageRelayHandler, CloudMessageRelayHandler itf.MessageRelayHandler) *Connector {
	return &Connector{
		FogMQTTClient: fogMqttClient,
		PlatformMQTTClient: platformMqttClient,
		PublishResultsToPlatform: publishResultsToPlatform,
		UserID: userID,
		LocalMessageRelayHandler: LocalMessageRelayHandler,
		CloudMessageRelayHandler: CloudMessageRelayHandler,
		SubscriptedLocalTopics: map[string]struct{}{},
	}
}

func (h *Connector) HandleOnDisconnect() {
	h.mu.Lock()
	defer h.mu.Unlock()
	clear(h.SubscriptedLocalTopics)
	logging.Logger.Debug("Cleared list of subscriptions to local topics")
}