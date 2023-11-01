package relay

func (relay *RelayController) processUserControlCommand(message []byte) {
	relay.Connector.ForwardToMaster(message)
}
