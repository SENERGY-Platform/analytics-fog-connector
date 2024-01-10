package mqtt 

import (
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/operator"
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/upstream"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"errors"	
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/logging"
)

type ReconnectHandler struct {
	UserID string 
}

func NewReconnectHandler(userID string) ReconnectHandler {
	return ReconnectHandler{
		UserID: userID,
	}
}

func (handler ReconnectHandler) Publish(client MQTT.Client, topic string) error {
	if !client.IsConnected() {
		logging.Logger.Error("WARNING: mqtt client not connected")
		return errors.New("mqtt client not connected")
	}

	token := client.Publish(topic, byte(2), false, "{}")
	if token.Wait() && token.Error() != nil {
		logging.Logger.Errorf("Error on Publish: ", token.Error())
		return token.Error()
	}
	logging.Logger.Debugf("Publish was successful")
	return nil
}

func (handler ReconnectHandler) RequestOperatorStatesSync(client MQTT.Client) error {
	logging.Logger.Debugf("Request operator sync")
	return handler.Publish(client, operator.GetOperatorControlSyncTriggerPubTopic(handler.UserID))
}

func (handler ReconnectHandler) RequestUpstreamForwardSync(client MQTT.Client) error {
	logging.Logger.Debugf("Request upstream forward sync")
	return handler.Publish(client, upstream.GetUpstreamControlSyncTriggerPubTopic(handler.UserID))
}

func (handler ReconnectHandler) OnConnect(client MQTT.Client) {
	logging.Logger.Debug("Connector connected to platform broker!")
	handler.RequestOperatorStatesSync(client)
	handler.RequestUpstreamForwardSync(client)
}

func OnConnectFog(client MQTT.Client) {
	logging.Logger.Debug("Connector connected to fog broker!")
}