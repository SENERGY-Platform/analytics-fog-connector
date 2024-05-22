package lib

import (
	"context"
	"errors"
	"io"
	"log"
	"syscall"

	"github.com/SENERGY-Platform/analytics-fog-connector/lib/clients/auth"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/clients/mqtt"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/config"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/connector"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/logging"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/relay"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/send_relay"
	mqttLib "github.com/SENERGY-Platform/analytics-fog-lib/lib/mqtt"
	"github.com/SENERGY-Platform/go-service-base/watchdog"
)

func Run(
	ctx    context.Context,
	stdout, stderr io.Writer,
	authClient auth.AuthClient,
	config config.Config,
) error {
	log.Println("Init Logger")
	err := logging.InitLogger(stdout, config.Debug)
	if err != nil {
		return err
	}

	watchdog := watchdog.New(syscall.SIGINT, syscall.SIGTERM)

	logging.Logger.Info("Get User ID")

	userID, err := authClient.GetUserID(config.Username, config.Password)
	if err != nil {
		logging.Logger.Error("Cant login user and get user id")
		return err
	}

	logging.Logger.Info("Setup Fog Broker connection")
	fogbrokerConfig := mqttLib.BrokerConfig(config.FogBroker)
	fogMqttClient := mqtt.NewFogMQTTClient(fogbrokerConfig, logging.Logger)
	watchdog.RegisterStopFunc(func() error {
		fogMqttClient.CloseConnection()
		return nil
	})

	logging.Logger.Info("Setup Platform Broker connection")
	PlatformBrokerConfig := mqttLib.BrokerConfig(config.PlatformBroker)
	platformMqttClient := mqtt.NewPlatformMQTTClient(PlatformBrokerConfig, userID, logging.Logger)
	watchdog.RegisterStopFunc(func() error {
		platformMqttClient.CloseConnection()
		return nil
	})

	localMqttClientPubF := func(topic string, data []byte) error {
		return fogMqttClient.Publish(topic, string(data), 2)
	}
	localMessageRelayHandler := send_relay.New(10000, localMqttClientPubF)
	watchdog.RegisterStopFunc(func() error {
		localMessageRelayHandler.Stop()
		return nil
	})

	cloudMqttClientPubF := func(topic string, data []byte) error {
		return platformMqttClient.Publish(topic, string(data), 2)
	}
	cloudMessageRelayHandler := send_relay.New(10000, cloudMqttClientPubF)
	watchdog.RegisterStopFunc(func() error {
		cloudMessageRelayHandler.Stop()
		return nil
	})

	logging.Logger.Info("Setup Connector, Upstream, Sync and Relay Controller")
	connector := connector.NewConnector(fogMqttClient, platformMqttClient, config.PublishResultsToPlatform, userID, localMessageRelayHandler, cloudMessageRelayHandler)
	relayController := relay.NewRelayController(connector, userID, config.PublishResultsToPlatform)

	fogMqttClient.SetRelayController(relayController)
	platformMqttClient.SetRelayController(relayController)

	logging.Logger.Info("Connect to brokers")
	fogMqttClient.ConnectMQTTBroker(nil, nil)
	platformMqttClient.ConnectMQTTBroker(&config.Username, &config.Password)

	cloudMessageRelayHandler.Start()
	localMessageRelayHandler.Start()

	logging.Logger.Info("Connector is ready")
	watchdog.Start()

	ec := watchdog.Join()
	logging.Logger.Info("Shutdowned graceful")
	if ec != 0 {
		return errors.New("Could not join")
	}
	return nil
}