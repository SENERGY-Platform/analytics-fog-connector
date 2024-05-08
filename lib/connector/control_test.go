package connector

import (
	"testing"

	"errors"

	"github.com/SENERGY-Platform/analytics-fog-connector/lib/clients/mqtt"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/send_relay"
	"github.com/stretchr/testify/assert"
)

func TestSendError(t *testing.T) {
	mqttClient := mqtt.MockMqttClient{
		ReturnError: true,
		Error: errors.New("cant connect"),
	}
	localMqttClientPubF := func(topic string, data []byte) error {
		return mqttClient.Publish(topic, string(data), 1)
	}
	localMessageRelayHandler := send_relay.New(10000, localMqttClientPubF)

	con := NewConnector(mqttClient, nil, false, "", localMessageRelayHandler, nil)

	err := con.ForwardStartOperatorToMaster([]byte(""))
	// when there is a mqtt error, we expect the relay handler to drop it and continue working normally
	assert.Nil(t, err)
}