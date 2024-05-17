package relay

import (
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/upstream"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/logging"
	"encoding/json"
)

func (relay *RelayController) processUpstreamDisable(message []byte) {
	logging.Logger.Debug("Received upstream disable command: " + string(message))

	disableMessage := upstream.UpstreamControlMessage{}
	err := json.Unmarshal(message, &disableMessage)
	if err != nil {
		logging.Logger.Error("Cant unmarshal disable upstream message: " + err.Error())
	}

	_ = relay.Connector.DisableForwarding(disableMessage)
}

func (relay *RelayController) processUpstreamEnable(message []byte) {
	logging.Logger.Debug("Received upstream enable command: " + string(message))

	enableMessage := upstream.UpstreamControlMessage{}
	err := json.Unmarshal(message, &enableMessage)
	if err != nil {
		logging.Logger.Error("Cant unmarshal enable upstream message: " + err.Error())
	}
	_ = relay.Connector.EnableForwarding(enableMessage)
}

func (relay *RelayController) processMessageToUpstream(message []byte, topic string) {
	logging.Logger.Debug("Received upstream message: " + string(message))
	_ = relay.Connector.ForwardOperatorResult(message, topic)
}

func (relay *RelayController) processUpstreamSync(message []byte) {
	logging.Logger.Debug("Received upstream sync response: " + string(message))

	syncMessage := upstream.UpstreamSyncMessage{}
	err := json.Unmarshal(message, &syncMessage)
	if err != nil {
		logging.Logger.Error("Cant unmarshal upstream sync message: " + err.Error())
	}
	err = relay.Connector.SyncUpstreamForward(syncMessage)
	if err != nil {
		logging.Logger.Error("Cant sync forwarding: " + err.Error())
	}
}
