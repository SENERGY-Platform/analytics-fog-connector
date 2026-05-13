package lib

import (
	"context"
	"errors"
	"io"
	"log"
	"syscall"

	"time"

	"github.com/SENERGY-Platform/analytics-fog-connector/internal/config"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/clients/auth"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/clients/mqtt"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/connector"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/logging"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/send_relay"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/subscription_handler"
	mqttLib "github.com/SENERGY-Platform/analytics-fog-lib/lib/mqtt"
	"github.com/SENERGY-Platform/go-service-base/watchdog"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func Run(
	ctx context.Context,
	stdout, stderr io.Writer,
	authClient auth.IDP,
	config config.Config,
) error {
	log.Println("Init Logger")
	err := logging.InitLogger(stdout, config.Debug)
	if err != nil {
		return err
	}

	watchdog := watchdog.New(syscall.SIGINT, syscall.SIGTERM)

	logging.Logger.Info("Get User ID")

	userID, err := authClient.GetUserID()
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
	subscriptionHandler := subscriptionhandler.NewSubscriptionHandler(connector, userID, config.PublishResultsToPlatform)
	// SubscriptionHandler for fog and platform each

	fogMqttClient.SetConnectionLostHandler(func(MQTT.Client) {
		logging.Logger.Debug("Connection to fog broker lost -> Clear cache of local operator topic subscriptions")
		connector.HandleOnDisconnect()
	})

	fogMqttClient.SetSubscriptionHandler(subscriptionHandler)

	platformMqttClient.SetSubscriptionHandler(subscriptionHandler)

	logging.Logger.Info("Connect to brokers")
	err = fogMqttClient.ConnectMQTTBroker(nil, nil)
	if err != nil {
		return err
	}
	if *config.PlatformBroker.Username == "" {
		//may be configured without credentials, if a proxy broker is used
		err = platformMqttClient.ConnectMQTTBroker(nil, nil)
		if err != nil {
			return err
		}
	} else {
		err = platformMqttClient.ConnectMQTTBroker(config.PlatformBroker.Username, config.PlatformBroker.Password)
		if err != nil {
			return err
		}
	}

	logging.Logger.Info("Start Periodic Sync of Operators and Upstream")
	go connector.RequestSyncs(ctx, time.Duration(config.SyncIntervalInSeconds))

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
