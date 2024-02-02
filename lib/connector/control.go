package connector

import (
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/operator"
)

func (connector *Connector) ForwardStartOperatorToMaster(message []byte) error {
	return connector.FogMQTTClient.Publish(operator.StartOperatorFogTopic, string(message), 2)
}

func (connector *Connector) ForwardStopOperatorToMaster(message []byte) error {
	return connector.FogMQTTClient.Publish(operator.StopOperatorFogTopic, string(message), 2)
}

func (connector *Connector) SyncOperatorStates(message []byte) error {
	return connector.FogMQTTClient.Publish(operator.OperatorControlSyncResponseFogTopic, string(message), 2)
}