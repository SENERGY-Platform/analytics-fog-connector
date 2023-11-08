package connector

import (
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/operator"
)

func (controller *Connector) ForwardToPlatform(message []byte) {
	controller.PlatformMQTTClient.Publish(operator.OperatorsResultTopic, string(message), 2)
}
