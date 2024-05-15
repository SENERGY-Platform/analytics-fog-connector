package relay

import (
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/logging"
)

func (relay *RelayController) processStartOperatorCommand(message []byte) {
	logging.Logger.Debug("Received operator start command: " + string(message))
	_ = relay.Connector.ForwardStartOperatorToMaster(message)
}

func (relay *RelayController) processStopOperatorCommand(message []byte) {
	logging.Logger.Debug("Received operator stop command: " + string(message))
	_ = relay.Connector.ForwardStopOperatorToMaster(message)
}

func (relay *RelayController) processOperatorSync(message []byte) {
	logging.Logger.Debug("Received operator sync response: " + string(message))
	_ = relay.Connector.SyncOperatorStates(message)
}
