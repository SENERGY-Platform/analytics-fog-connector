package e2e

import (
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/connector"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/relay"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/send_relay"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/logging"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/clients/mqtt"
	"github.com/SENERGY-Platform/analytics-fog-connector/tests/dependencies"
	mqttLib "github.com/SENERGY-Platform/analytics-fog-lib/lib/mqtt"
	srv_base "github.com/SENERGY-Platform/go-service-base/util"
	"context"
	"fmt"
	"testing"
)

type Env struct {
	fogMqttClient mqtt.MQTTClient
	platformMqttClient mqtt.MQTTClient
	vernemq *dependencies.VerneMQ
	mosquitto *dependencies.Mosquitto
	UserID string
}

func NewEnv(ctx context.Context, t *testing.T) (*Env, error) {
	_, err := logging.InitLogger(srv_base.LoggerConfig{})
	if err != nil {
		return &Env{}, err
	}

	vernemq, err := dependencies.NewVerneMQ(ctx)
	if err != nil {
		return &Env{}, err
	}
	mosquitto, err := dependencies.NewMosquitto(ctx)
	if err != nil {
		return &Env{}, err
	}

	return &Env{
		vernemq: vernemq,
		mosquitto: mosquitto,
		UserID: "user",
	}, nil
}

func (e *Env) Start(ctx context.Context, t *testing.T) (error, string, string) {
	publishResultsToPlatform := true 

	fmt.Println("Start Verne")
	err, vernePort := e.vernemq.StartAndWait(ctx)
	if err != nil {
		return err, "", ""
	}
	
	fmt.Println("Start Mosquitto")
	err, mosquittoPort := e.mosquitto.StartAndWait(ctx)
	if err != nil {
		return err, "", ""
	}
	
	fogbrokerConfig := mqttLib.BrokerConfig(mqttLib.FogBrokerConfig{
		Host: "localhost",
		Port: mosquittoPort,
	})
	fogMqttClient := mqtt.NewFogMQTTClient(fogbrokerConfig, logging.Logger)
	localMqttClientPubF := func(topic string, data []byte) error {
		return fogMqttClient.Publish(topic, string(data), 1)
	}
	localMessageRelayHandler := send_relay.New(10000, localMqttClientPubF)


	PlatformBrokerConfig := mqttLib.BrokerConfig(mqttLib.PlatformBrokerConfig{
		Host: "localhost",
		Port: vernePort,
	})
	platformMqttClient := mqtt.NewPlatformMQTTClient(PlatformBrokerConfig, e.UserID, logging.Logger)
	cloudMqttClientPubF := func(topic string, data []byte) error {
		return platformMqttClient.Publish(topic, string(data), 1)
	}
	cloudMessageRelayHandler := send_relay.New(10000, cloudMqttClientPubF)
	
	connector := connector.NewConnector(fogMqttClient, platformMqttClient, publishResultsToPlatform, e.UserID, localMessageRelayHandler, cloudMessageRelayHandler)
	relayController := relay.NewRelayController(connector, e.UserID, publishResultsToPlatform)

	fogMqttClient.SetRelayController(relayController)
	platformMqttClient.SetRelayController(relayController)

	fmt.Println("Connect to local broker")
	fogMqttClient.ConnectMQTTBroker(nil, nil)
	fmt.Println("Connect to cloud broker")
	platformMqttClient.ConnectMQTTBroker(nil, nil)
	return nil, vernePort, mosquittoPort
}

func (e *Env) PublishToCloud(topic string, payload []byte) {
	e.platformMqttClient.Publish(topic, string(payload), 2)
}

func (e *Env) PublishToFog(topic string, payload []byte) {
	e.fogMqttClient.Publish(topic, string(payload), 2)

}