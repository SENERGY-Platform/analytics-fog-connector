package relay

func (relay *RelayController) processOperatorControlCommand(message []byte) {
	relay.Connector.ForwardToMaster(message)
}

func (relay *RelayController) processOperatorOutputMessage(message []byte) {
	relay.Connector.ForwardToPlatform(message)
}