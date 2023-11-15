package relay 

import (
	downstreamLib "github.com/SENERGY-Platform/analytics-fog-lib/lib/downstream"
	"encoding/json"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/logging"

)

func (relay *RelayController) processDownstreamMessage(message []byte) {
	downStreamEnvelope := downstreamLib.DownstreamEnvelope{}
	err := json.Unmarshal(message, &downStreamEnvelope)
	if err != nil {
		logging.Logger.Errorf("Cant unmarshal enable downstream message:", err)
	}
	_ = relay.Connector.ForwardCloudMessageToFog(downStreamEnvelope)
}