package relay

import (
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/upstream"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/logging"
	"encoding/json"
)

func (relay *RelayController) processUpstreamDisable(message []byte) {
	disableMessage := upstream.UpstreamControlMessage{}
	err := json.Unmarshal(message, &disableMessage)
	if err != nil {
		logging.Logger.Errorf("Cant unmarshal disable upstream message:", err)
	}

	_ = relay.Connector.DisableForwarding(disableMessage)
}

func (relay *RelayController) processUpstreamEnable(message []byte) {
	enableMessage := upstream.UpstreamControlMessage{}
	err := json.Unmarshal(message, &enableMessage)
	if err != nil {
		logging.Logger.Errorf("Cant unmarshal enable upstream message:", err)
	}
	_ = relay.Connector.EnableForwarding(enableMessage)
}

func (relay *RelayController) processMessageToUpstream(message []byte, topic string) {
	_ = relay.Connector.ForwardOperatorResult(message, topic)
}

func (relay *RelayController) processUpstreamSync(message []byte) {
	syncMessage := upstream.UpstreamSyncMessage{}
	err := json.Unmarshal(message, &syncMessage)
	if err != nil {
		logging.Logger.Errorf("Cant unmarshal upstream sync message:", err)
	}
	_ = relay.Connector.SyncUpstreamForward(syncMessage)
}
