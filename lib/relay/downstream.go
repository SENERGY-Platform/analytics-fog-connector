package relay 

import (
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/logging"
)

func (relay *RelayController) processOperatorDownstreamMessage(message []byte, topic string) {
	downStreamMessage := string(message)
	err := relay.Connector.ForwardCloudMessageToFog(downStreamMessage, topic)
	if err != nil {
		logging.Logger.Errorf("Cant forward cloud message to fog: ", err.Error())
	}
}