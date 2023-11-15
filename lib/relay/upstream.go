package relay

import (
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/upstream"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/logging"
	"encoding/json"
)

func (relay *RelayController) processUpstreamDisable(message []byte) {
	disableMessage := upstream.UpstreamDisableMessage{}
	err := json.Unmarshal(message, &disableMessage)
	if err != nil {
		logging.Logger.Errorf("Cant unmarshal disable upstream message:", err)
	}

	_ = relay.Connector.DisableForwarding(disableMessage)
}

func (relay *RelayController) processUpstreamEnable(message []byte) {
	enableMessage := upstream.UpstreamEnableMessage{}
	err := json.Unmarshal(message, &enableMessage)
	if err != nil {
		logging.Logger.Errorf("Cant unmarshal enable upstream message:", err)
	}
	_ = relay.Connector.EnableForwarding(enableMessage)
}

func (relay *RelayController) processMessageToUpstream(message []byte, topic string) {
	_ = relay.Connector.ForwardOperatorResult(message, topic)
}