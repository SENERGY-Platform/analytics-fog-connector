package connector

import (
	operator "github.com/SENERGY-Platform/analytics-fog-lib/lib/operator"
	deployLocationLib "github.com/SENERGY-Platform/analytics-fog-lib/lib/location"
	"strings"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/logging"
)


func (connector *Connector) ForwardCloudMessageToFog(downstreamMessage string, topic string) error {
	logging.Logger.Debugf("Received cloud message to be forwarded to fog broker at topic: ", topic)
	baseOperatorName, operatorID, pipelineID := GetOperatorIDsFromTopic(topic)
	fogTopic, err := operator.GenerateOperatorOutputTopic(baseOperatorName, "", operatorID, deployLocationLib.Local)
	if err != nil {
		return err 
	}

	fogTopic = fogTopic + "/" + pipelineID
	// pipeline ID check 

	logging.Logger.Debugf("Try to publish downstream message: %s to fog topic: %s", downstreamMessage, fogTopic)
	err = connector.FogMQTTClient.Publish(fogTopic, downstreamMessage, 2)
	return err 
}

func GetOperatorIDsFromTopic(topic string) (string, string, string) {
	split := strings.Split(topic, "/")
	operatorName := split[len(split)-3]
	operatorID := split[len(split)-2]
	pipelineID := split[len(split)-1]
	return operatorName, operatorID, pipelineID
}