package connector

import (
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/operator"
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/upstream"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/logging"
	"strings"
	"encoding/json"
)

func (connector *Connector) ForwardOperatorResult(message []byte, topic string) error {
	// Wrap operator result in an envelope containing the required ID to map on cloud side
	// And publish to senergy-connector or directly the cloud broker
	// On Cloud side, the operator name and the operator ID are required to correctly distribute the messages
	envelope, err := connector.CreateUpstreamEnvelope(message, topic)
	if err != nil {
		logging.Logger.Errorf("Can create upstream envelope: %s", err)
		return err
	}
	envelopeMessage, err := json.Marshal(envelope)
	if err != nil {
		logging.Logger.Errorf("Cant marshal upstream envelope: %s", err)
		return err
	}
	envelopeString := string(envelopeMessage)

	// Topic which collects all mesagess across operators that must be forwarded by the senergy-connector to the platform
	logging.Logger.Debugf("Try to publish upstream message to fog broker: %s", envelopeString)
	err = connector.FogMQTTClient.Publish(operator.OperatorsResultTopic, envelopeString, 2)

	if err != nil {
		logging.Logger.Errorf("Cant publish upstream message to fog broker")
		return err
	}

	if connector.PublishResultsToPlatform {
		logging.Logger.Debugf("Try to publish upstream message to platform: %s", envelopeString)
		err = connector.PlatformMQTTClient.Publish(operator.OperatorsResultTopic, envelopeString, 2)
		if err != nil {
			logging.Logger.Errorf("Cant publish upstream message to platform broker")
			return err
		}
	}

	return nil
}

func GetOperatorNameAndIDFromTopic(topic string) (string, string) {
	split := strings.Split(topic, "/")
	operatorID := split[len(split)-2]
	operatorName := split[len(split)-3]
	return operatorName, operatorID
}

func (connector *Connector) CreateUpstreamEnvelope(message []byte, topic string) (upstream.UpstreamEnvelope, error) {
	baseOperatorName, operatorID := GetOperatorNameAndIDFromTopic(topic)
	return upstream.UpstreamEnvelope{
		Message: string(message),
		BaseOperatorName: baseOperatorName,
		OperatorID: operatorID,
	}, nil
}


// TODO !!! store topics on disk that must be forwarded
func (connector *Connector) EnableForwarding(enableMessage upstream.UpstreamEnableMessage) error {
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

func (connector *Connector) DisableForwarding(disableMessage upstream.UpstreamDisableMessage) error {
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