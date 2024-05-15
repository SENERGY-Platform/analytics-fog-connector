package mqtt

import (
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/operator"
	MQTT "github.com/eclipse/paho.mqtt.golang"

	"github.com/SENERGY-Platform/analytics-fog-lib/lib/mqtt"
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/upstream"
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/downstream"
	"log/slog"
)

func NewMQTTClient(brokerConfig mqtt.BrokerConfig, logger *slog.Logger, topicConfig mqtt.TopicConfig, onConnectHandler func(MQTT.Client), subscribeInitial bool) *mqtt.MQTTClient {
	return &mqtt.MQTTClient{
		Broker:      brokerConfig,
		TopicConfig: topicConfig,
		Logger:      logger,
		OnConnectHandler: onConnectHandler,
		SubscribeInitial: subscribeInitial,
	}
}

func NewPlatformMQTTClient(brokerConfig mqtt.BrokerConfig, userID string, logger *slog.Logger) *mqtt.MQTTClient {
	topicConfig := mqtt.TopicConfig{
		// Start/Stop operator commands
		operator.GetStartOperatorCloudTopic(userID): byte(2),
		operator.GetStopOperatorCloudTopic(userID): byte(2),

		// Cloud operator output that are need in fog  
		downstream.GetDownstreamOperatorCloudSubTopic(userID): byte(2),

		// Commands to enable and disbale forwarding of specific fog operators that are needed in cloud
		upstream.GetUpstreamEnableCloudTopic(userID): byte(2),
		upstream.GetUpstreamDisableCloudTopic(userID): byte(2),

		// Upstream Sync Response
		upstream.GetUpstreamControlSyncResponseTopic(userID): byte(2),

		// Operator Sync Response
		operator.GetOperatorControlSyncResponseTopic(userID): byte(2),
	}

	onConnectHandler := NewReconnectHandler(userID)
	return NewMQTTClient(brokerConfig, logger, topicConfig, onConnectHandler.OnConnect, true)
}

func NewFogMQTTClient(brokerConfig mqtt.BrokerConfig, logger *slog.Logger) *mqtt.MQTTClient {
	topicConfig := mqtt.TopicConfig{
	}
	return NewMQTTClient(brokerConfig, logger, topicConfig, OnConnectFog, false)
}