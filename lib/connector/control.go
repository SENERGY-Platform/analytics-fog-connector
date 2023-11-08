package connector

import "github.com/SENERGY-Platform/analytics-fog-lib/lib/operator"

func (controller *Connector) ForwardToMaster(message []byte) {
	controller.FogMQTTClient.Publish(operator.OperatorsControlTopic, string(message), 2)
}
