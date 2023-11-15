package connector

import (
	downstreamLib "github.com/SENERGY-Platform/analytics-fog-lib/lib/downstream"
	topicLib "github.com/SENERGY-Platform/analytics-fog-lib/lib/topic"
	deployLocationLib "github.com/SENERGY-Platform/analytics-fog-lib/lib/location"
)


func (connector *Connector) ForwardCloudMessageToFog(downstreamMessage downstreamLib.DownstreamEnvelope) error {
	topic := topicLib.GenerateOperatorOutputTopic(downstreamMessage.BaseOperatorName, downstreamMessage.OperatorID, deployLocationLib.Local)

	err := connector.FogMQTTClient.Publish(topic, downstreamMessage.Message, 2)
	return err 
}