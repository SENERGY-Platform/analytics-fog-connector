package subscriptionhandler

import (
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/logging"
)

func (handler *SubscriptionHandler) processStartOperatorCommand(message []byte) {
	logging.Logger.Debug("Received operator start command: " + string(message))
	_ = handler.Connector.ForwardStartOperatorToMaster(message)
}

func (handler *SubscriptionHandler) processStopOperatorCommand(message []byte) {
	logging.Logger.Debug("Received operator stop command: " + string(message))
	_ = handler.Connector.ForwardStopOperatorToMaster(message)
}

func (handler *SubscriptionHandler) processOperatorSync(message []byte) {
	logging.Logger.Debug("Received operator sync response: " + string(message))
	_ = handler.Connector.SyncOperatorStates(message)
}
