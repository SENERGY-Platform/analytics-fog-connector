package connector

import (
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/operator"

)

func (connector *Connector) ForwardStartOperatorToMaster(payload []byte) error {
	message := Message{
		topic: operator.StartOperatorFogTopic,
		payload: payload,
	}
	return connector.LocalMessageRelayHandler.Put(message)
}

func (connector *Connector) ForwardStopOperatorToMaster(payload []byte) error {
	message := Message{
		topic: operator.StopOperatorFogTopic,
		payload: payload,
	}
	return connector.LocalMessageRelayHandler.Put(message)
}

func (connector *Connector) SyncOperatorStates(payload []byte) error {
	message := Message{
		topic: operator.OperatorControlSyncResponseFogTopic,
		payload: payload,
	}
	return connector.LocalMessageRelayHandler.Put(message)
}