package connector

import (
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/upstream"
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/location"
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/operator"

	"github.com/SENERGY-Platform/analytics-fog-connector/lib/logging"
	"strings"
)

func (connector *Connector) ForwardOperatorResult(payload []byte, fogTopic string) error {
	baseOperatorName := GetOperatorNameFromTopic(fogTopic)
	message := string(payload)
	cloudOperatorTopic, _ := operator.GenerateOperatorOutputTopic(baseOperatorName, "", "", location.Cloud)
	platformTopic := upstream.CloudUpstreamTopic + "/" + cloudOperatorTopic
	logging.Logger.Debugf("Try to publish upstream message: %s to platform topic: %s", message, platformTopic)
	err := connector.PlatformMQTTClient.Publish(platformTopic, message, 2)
	if err != nil {
		logging.Logger.Errorf("Cant publish upstream message to platform broker ")
		return err
	}
	logging.Logger.Debugf("Successfully published upstream message")
	return nil
}

func GetOperatorNameFromTopic(topic string) (string) {
	split := strings.Split(topic, "/")
	operatorName := split[len(split)-3]
	return operatorName
}

// TODO !!! store topics on disk that must be forwarded
// or on startup get current list with topics
func (connector *Connector) EnableForwarding(enableMessage upstream.UpstreamControlMessage) error {
	topic := enableMessage.OperatorOutputTopic
	logging.Logger.Infof("Try to subscribe to %s", topic)
	err := connector.FogMQTTClient.Subscribe(topic, 2)
	if err != nil {
		logging.Logger.Errorf("Cant subscribe to operator output topic %s: %s", topic, err)
		return err
	}
	logging.Logger.Infof("Successfully subscribed to %s", topic)

	return nil
}

func (connector *Connector) DisableForwarding(disableMessage upstream.UpstreamControlMessage) error {
	topic := disableMessage.OperatorOutputTopic
	logging.Logger.Infof("Try to unsubscribe from %s", topic)

	err := connector.FogMQTTClient.Unsubscribe(topic)
	if err != nil {
		logging.Logger.Errorf("Cant unsubscribe from operator output topic %s: %s", topic, err)
		return err
	}
	logging.Logger.Infof("Successfully unsubscribed from %s", topic)
	return nil
}

// TODO unsubscribe 
func (connector *Connector) SyncUpstreamForward(syncMessage upstream.UpstreamSyncMessage) error {
	for _, topic := range(syncMessage.OperatorOutputTopics) {
		err := connector.EnableForwarding(upstream.UpstreamControlMessage{
			OperatorOutputTopic: topic,
		})
		if err != nil {
			return err
		}
	}
	return nil
}


