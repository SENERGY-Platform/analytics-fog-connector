package connector

import (
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/operator"
	"encoding/json"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/logging"
)

func (connector *Connector) ForwardStartOperatorToMaster(message []byte) error {
	return connector.FogMQTTClient.Publish(operator.StartOperatorFogTopic, string(message), 2)
}

func (connector *Connector) ForwardStopOperatorToMaster(message []byte) error {
	return connector.FogMQTTClient.Publish(operator.StopOperatorFogTopic, string(message), 2)
}

func (connector *Connector) SyncOperatorStates(syncMsg []operator.StartOperatorControlCommand) error {
	for _, operatorStartCmd := range(syncMsg) {
		message, err := json.Marshal(operatorStartCmd)
		if err != nil {
			logging.Logger.Errorf("Cant marshal sync operator start message: " + err.Error())
			return err
		}
		connector.ForwardStartOperatorToMaster(message)
	}
	return nil
}