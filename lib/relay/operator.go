package relay

import (

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
