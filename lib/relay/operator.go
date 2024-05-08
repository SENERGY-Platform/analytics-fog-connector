package relay

import (
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/logging"
)

func (relay *RelayController) processStartOperatorCommand(message []byte) {
	logging.Logger.Debugf("Received operator start command: %s", string(message))
	_ = relay.Connector.ForwardStartOperatorToMaster(message)
}

func (relay *RelayController) processStopOperatorCommand(message []byte) {
	logging.Logger.Debugf("Received operator stop command: %s", string(message))
	_ = relay.Connector.ForwardStopOperatorToMaster(message)
}

func (relay *RelayController) processOperatorSync(message []byte) {
	logging.Logger.Debugf("Received operator sync response: %s", string(message))
	_ = relay.Connector.SyncOperatorStates(message)
}
