package relay

func (relay *RelayController) processStartOperatorCommand(message []byte) {
	_ = relay.Connector.ForwardStartOperatorToMaster(message)
}

func (relay *RelayController) processStopOperatorCommand(message []byte) {
	_ = relay.Connector.ForwardStopOperatorToMaster(message)
}
