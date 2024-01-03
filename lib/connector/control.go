package connector

import (
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/operator"

)

func (controller *Connector) ForwardStartOperatorToMaster(message []byte) error {
	return controller.FogMQTTClient.Publish(operator.StartOperatorFogTopic, string(message), 2)
}

func (controller *Connector) ForwardStopOperatorToMaster(message []byte) error {
	return controller.FogMQTTClient.Publish(operator.StopOperatorFogTopic, string(message), 2)
}
