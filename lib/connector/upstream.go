package connector

import (
	"fmt"

	"github.com/SENERGY-Platform/analytics-fog-lib/lib/location"
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/operator"
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/upstream"

	"strings"

	"github.com/SENERGY-Platform/analytics-fog-connector/lib/logging"
)

func (connector *Connector) ForwardOperatorResult(payload []byte, fogTopic string) error {
	baseOperatorName := GetOperatorNameFromTopic(fogTopic)
	cloudOperatorTopic, _ := operator.GenerateOperatorOutputTopic(baseOperatorName, "", "", location.Cloud)
	platformTopic := upstream.CloudUpstreamTopic + "/" + cloudOperatorTopic
	logging.Logger.Debug(fmt.Sprintf("Try to publish upstream message: %s to platform topic: %s", string(payload), platformTopic))
	message := Message{
		topic: platformTopic,
		payload: payload,
	}
	err := connector.CloudMessageRelayHandler.Put(message)
	if err != nil {
		logging.Logger.Error("Cant publish upstream message to platform broker: " + err.Error())
		return err
	}
	logging.Logger.Debug("Successfully published upstream message")
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
	logging.Logger.Info("Try to subscribe to: " + topic)
	err := connector.FogMQTTClient.Subscribe(topic, 2)
	if err != nil {
		logging.Logger.Error(fmt.Sprintf("Cant subscribe to operator output topic %s: %s", topic, err))
		return err
	}
	logging.Logger.Info("Successfully subscribed to: " + topic)

	return nil
}

func (connector *Connector) DisableForwarding(disableMessage upstream.UpstreamControlMessage) error {
	topic := disableMessage.OperatorOutputTopic
	logging.Logger.Info("Try to unsubscribe from: " + topic)

	err := connector.FogMQTTClient.Unsubscribe(topic)
	if err != nil {
		logging.Logger.Error(fmt.Sprintf("Cant unsubscribe from operator output topic %s: %s", topic, err))
		return err
	}
	logging.Logger.Info("Successfully unsubscribed from: " + topic)
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


