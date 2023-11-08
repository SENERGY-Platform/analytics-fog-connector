package mqtt

import (
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/control"
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/operator"

	"github.com/SENERGY-Platform/analytics-fog-lib/lib/mqtt"
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
	userControlTopic := control.GetConnectorControlTopic(userID)
	topicConfig := mqtt.TopicConfig{
		userControlTopic: byte(2),
	}

	return NewMQTTClient(brokerConfig, logger, topicConfig)
}

func NewFogMQTTClient(brokerConfig mqtt.BrokerConfig, logger *log_level.Logger) *mqtt.MQTTClient {
	topicConfig := mqtt.TopicConfig{
		operator.OperatorsResultTopic: byte(2),
	}
	return NewMQTTClient(brokerConfig, logger, topicConfig)
}