package subscriptionhandler 

import (
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/logging"
)

func (handler *SubscriptionHandler) processOperatorDownstreamMessage(message []byte, topic string) {
	logging.Logger.Debug("Received operator downstream command: " + string(message))
	downStreamMessage := string(message)
	err := handler.Connector.ForwardCloudMessageToFog(downStreamMessage, topic)
	if err != nil {
		logging.Logger.Error("Cant forward cloud message to fog: " + err.Error())
	}
}