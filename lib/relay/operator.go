package relay

func (relay *RelayController) processOperatorControlCommand(message []byte) {
	_ = relay.Connector.ForwardToMaster(message)
}

