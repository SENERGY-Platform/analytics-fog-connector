package connector

import (
	"fmt"

	"github.com/SENERGY-Platform/analytics-fog-lib/lib/operator"
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/upstream"

	"strings"

	"github.com/SENERGY-Platform/analytics-fog-connector/lib/logging"
)

func (connector *Connector) ForwardOperatorResult(payload []byte, fogTopic string) error {
	baseOperatorName := GetOperatorNameFromTopic(fogTopic)
	cloudOperatorTopic := operator.GenerateCloudOperatorTopic(baseOperatorName)
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

func (connector *Connector) EnableForwarding(enableMessage upstream.UpstreamControlMessage) error {
	connector.mu.Lock()
	defer connector.mu.Unlock()
	topic := enableMessage.OperatorOutputTopic
	logging.Logger.Info("Try to subscribe to: " + topic)
	qos := 2
	err := connector.FogMQTTClient.Subscribe(topic, qos)
	if err != nil {
		logging.Logger.Error(fmt.Sprintf("Cant subscribe to operator output topic %s: %s", topic, err))
		return err
	}
	connector.SubscriptedLocalTopics[topic] = struct{}{}
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

func (connector *Connector) SyncUpstreamForward(syncMessage upstream.UpstreamSyncMessage) error {
	logging.Logger.Debug("Sync upstream forwardings")
	topicsWithMissingForwarding := []string{}
	connector.mu.Lock()
	currentSubscriptedTopics := map[string]struct{}{} 
	for topic, _ := range(connector.SubscriptedLocalTopics) {
		currentSubscriptedTopics[topic] = struct{}{}
	}
	connector.mu.Unlock()
	
	for _, expectedTopic := range(syncMessage.OperatorOutputTopics) {
		_, ok := currentSubscriptedTopics[expectedTopic]; if !ok {
			topicsWithMissingForwarding = append(topicsWithMissingForwarding, expectedTopic)
		}
	}
	connector.EnableMissingForwarding(topicsWithMissingForwarding)

	topicsWithOrphanForwardings := []string{}
	for topicWithActiveForwarding, _ := range(currentSubscriptedTopics) {
		topicShallBeForwarded := false
		for _, expectedTopic := range(syncMessage.OperatorOutputTopics) {
			if expectedTopic == topicWithActiveForwarding {
				topicShallBeForwarded = true
				break
			}
		}

		if !topicShallBeForwarded {
			topicsWithOrphanForwardings = append(topicsWithOrphanForwardings, topicWithActiveForwarding)
		}
	}
	connector.DisableOrphanForwarding(topicsWithOrphanForwardings)
	return nil
}

func (connector *Connector) EnableMissingForwarding(operatorTopics []string) error {
	logging.Logger.Debug("Enable missing upstream forwardings for: " + strings.Join(operatorTopics, ","))
	for _, topic := range(operatorTopics) {
		logging.Logger.Debug("Try to enable forwarding for: " + topic)
		err := connector.EnableForwarding(upstream.UpstreamControlMessage{
			OperatorOutputTopic: topic,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (connector *Connector) DisableOrphanForwarding(operatorTopics []string) error {
	logging.Logger.Debug("Disable orphan upstream forwardings for: " + strings.Join(operatorTopics, ","))
	for _, topic := range(operatorTopics) {
		logging.Logger.Debug("Try to disable forwarding for: " + topic)
		err := connector.DisableForwarding(upstream.UpstreamControlMessage{
			OperatorOutputTopic: topic,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (connector *Connector) RequestUpstreamSync() error {
	topic := upstream.GetUpstreamControlSyncTriggerPubTopic(connector.UserID)
	return connector.PlatformMQTTClient.Publish(topic, "", 2)
}