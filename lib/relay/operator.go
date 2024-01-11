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
	_ = relay.Connector.SyncOperatorStates(message)
}
