package subscriptionhandler

import (
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/upstream"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/logging"
	"encoding/json"
)

func (handler *SubscriptionHandler) processUpstreamDisable(message []byte) {
	logging.Logger.Debug("Received upstream disable command: " + string(message))

	disableMessage := upstream.UpstreamControlMessage{}
	err := json.Unmarshal(message, &disableMessage)
	if err != nil {
		logging.Logger.Error("Cant unmarshal disable upstream message: " + err.Error())
	}

	_ = handler.Connector.DisableForwarding(disableMessage)
}

func (handler *SubscriptionHandler) processUpstreamEnable(message []byte) {
	logging.Logger.Debug("Received upstream enable command: " + string(message))

	enableMessage := upstream.UpstreamControlMessage{}
	err := json.Unmarshal(message, &enableMessage)
	if err != nil {
		logging.Logger.Error("Cant unmarshal enable upstream message: " + err.Error())
	}
	_ = handler.Connector.EnableForwarding(enableMessage)
}

func (handler *SubscriptionHandler) processMessageToUpstream(message []byte, topic string) {
	logging.Logger.Debug("Received upstream message: " + string(message))
	_ = handler.Connector.ForwardOperatorResult(message, topic)
}

func (handler *SubscriptionHandler) processUpstreamSync(message []byte) {
	logging.Logger.Debug("Received upstream sync response: " + string(message))

	syncMessage := upstream.UpstreamSyncMessage{}
	err := json.Unmarshal(message, &syncMessage)
	if err != nil {
		logging.Logger.Error("Cant unmarshal upstream sync message: " + err.Error())
	}
	err = handler.Connector.SyncUpstreamForward(syncMessage)
	if err != nil {
		logging.Logger.Error("Cant sync forwarding: " + err.Error())
	}
}
