package mqtt

import (
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/control"

	"github.com/SENERGY-Platform/analytics-fog-lib/lib/mqtt"
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/upstream"
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/downstream"

	log_level "github.com/y-du/go-log-level"
)

func NewMQTTClient(brokerConfig mqtt.BrokerConfig, logger *log_level.Logger, topicConfig mqtt.TopicConfig) *mqtt.MQTTClient {
	return &mqtt.MQTTClient{
		Broker:      brokerConfig,
		TopicConfig: topicConfig,
		Logger:      logger,
	}
}

func NewPlatformMQTTClient(brokerConfig mqtt.BrokerConfig, userID string, logger *log_level.Logger) *mqtt.MQTTClient {
	topicConfig := mqtt.TopicConfig{
		// Start/Stop operator commands
		control.GetConnectorControlTopic(userID): byte(2),
		
		// Cloud operator/device output that are need in fog  
		downstream.GetDownstreamTopic(userID): byte(2),

		// Commands to enable and disbale forwarding of specific fog operators that are needed in cloud
		upstream.GetUpstreamProxyEnableTopic(userID): byte(2),
		upstream.GetUpstreamProxyDisableTopic(userID): byte(2),
	}

	return NewMQTTClient(brokerConfig, logger, topicConfig)
}

func NewFogMQTTClient(brokerConfig mqtt.BrokerConfig, logger *log_level.Logger) *mqtt.MQTTClient {
	topicConfig := mqtt.TopicConfig{
	}
	return NewMQTTClient(brokerConfig, logger, topicConfig)
}