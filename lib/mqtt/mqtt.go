package mqtt

import (
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/control"
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/mqtt"
	log_level "github.com/y-du/go-log-level"
)

func NewMQTTClient(brokerConfig mqtt.BrokerConfig, userID string, logger *log_level.Logger) *mqtt.MQTTClient {
	userControlTopic := control.GetConnectorControlTopic(userID)
	logger.Debug(userControlTopic)
	topics := mqtt.TopicConfig{
		userControlTopic: byte(2),
	}

	return &mqtt.MQTTClient{
		Broker:      brokerConfig,
		TopicConfig: topics,
		Logger:      logger,
	}
}
