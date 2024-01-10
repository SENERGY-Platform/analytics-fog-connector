package relay

import (
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/logging"
	"encoding/json"
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/operator"

)

func (relay *RelayController) processStartOperatorCommand(message []byte) {
	_ = relay.Connector.ForwardStartOperatorToMaster(message)
}

func (relay *RelayController) processStopOperatorCommand(message []byte) {
	_ = relay.Connector.ForwardStopOperatorToMaster(message)
}

func (relay *RelayController) processOperatorSync(message []byte) {
	syncMessage := []operator.StartOperatorControlCommand{}
	err := json.Unmarshal(message, &syncMessage)
	if err != nil {
		logging.Logger.Errorf("Cant unmarshal upstream sync message:", err)
	}
	_ = relay.Connector.SyncOperatorStates(syncMessage)
}
