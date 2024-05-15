package connector

import (
	"fmt"
	"strings"

	"github.com/SENERGY-Platform/analytics-fog-connector/lib/logging"
	deployLocationLib "github.com/SENERGY-Platform/analytics-fog-lib/lib/location"
	operator "github.com/SENERGY-Platform/analytics-fog-lib/lib/operator"
)


func (connector *Connector) ForwardCloudMessageToFog(payload string, topic string) error {
	logging.Logger.Debug("Received cloud message to be forwarded to fog broker at topic: " + topic)
	baseOperatorName, operatorID, pipelineID := GetOperatorIDsFromTopic(topic)
	fogTopic, err := operator.GenerateOperatorOutputTopic(baseOperatorName, "", operatorID, deployLocationLib.Local)
	if err != nil {
		return err 
	}

	fogTopic = fogTopic + "/" + pipelineID
	// pipeline ID check 

	message := Message{
		topic: fogTopic,
		payload: []byte(payload),
	}

	logging.Logger.Debug(fmt.Sprintf("Try to publish downstream message: %s to fog topic: %s", payload, fogTopic))
	err = connector.LocalMessageRelayHandler.Put(message)
	return err 
}

func GetOperatorIDsFromTopic(topic string) (string, string, string) {
	split := strings.Split(topic, "/")
	operatorName := split[len(split)-3]
	operatorID := split[len(split)-2]
	pipelineID := split[len(split)-1]
	return operatorName, operatorID, pipelineID
}