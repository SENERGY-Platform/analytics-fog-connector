package e2e

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
	"os"

	"github.com/SENERGY-Platform/analytics-fog-connector/lib"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/clients/auth"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/config"
	"github.com/SENERGY-Platform/analytics-fog-connector/tests/dependencies"
	mqttLib "github.com/SENERGY-Platform/analytics-fog-lib/lib/mqtt"
	testLib "github.com/hahahannes/e2e-go-utils/lib"
	"github.com/hahahannes/e2e-go-utils/lib/streaming/mqtt"
)

type ApplicationLogger struct {
	readyLogChannel chan string
	logChannel chan string
	customChannel chan string
}

func (l *ApplicationLogger) Write(p []byte) (n int, err error) {
	msg := string(p)
	select {
		// use select default to not block logging when readylogchannel has no receiver anymore after the check
		case l.readyLogChannel <- msg:
		default:
	}

	select {
		case l.logChannel <- msg:
		default:
	}

	select {
		case l.customChannel <- msg:
		default:
	}
	
	return 1, nil
}

func (e *Env) StartAndWait(ctx context.Context, t *testing.T, customChannel chan string) error {
	readyLogChannel := make(chan string, 100)
	logChannel := make(chan string)
	logger := &ApplicationLogger{
		readyLogChannel: readyLogChannel,
		logChannel: logChannel,
		customChannel: customChannel,
	}

	go func() {
		f, err := os.OpenFile("connector.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Print("Could not open log file")
		}
		
		for {
			select {
			case log := <- logChannel:
				if _, err := f.Write([]byte(log + "\n")); err != nil {
					fmt.Print("Could not write to log file")
				}
			case <- ctx.Done():
				if err := f.Close(); err != nil {
					fmt.Print("Could not close file")
				}
				return
			}
		}
		
	}()

	authMock := auth.MockAuth{
		UserID: e.UserID,
	}
	runConfig := config.Config{
		Debug: true,
		DataDir: ".",
		FogBroker: mqttLib.FogBrokerConfig{Host:"localhost", Port: e.fogBrokerPort},
		PlatformBroker: mqttLib.PlatformBrokerConfig{Host:"localhost", Port: e.cloudBrokerPort},
	}
	
	received, err := testLib.WaitForStringReceived(".*Connector is ready.*", func (sendCtx context.Context) error {
		return lib.Run(sendCtx, logger, logger, authMock, runConfig)
	}, readyLogChannel, 30 * time.Second, true)

	if err != nil {
		return err 
	}

	if received.Received == false {
		return errors.New("Connector ready log not received")
	}
	t.Log("Connector subscribed log received!")

	return nil
}

type Env struct {
	fogMqttClient *mqtt.MQTTClient
	platformMqttClient *mqtt.MQTTClient
	cloudBroker *dependencies.Mosquitto
	fogBroker *dependencies.Mosquitto
	UserID string
	cloudBrokerPort string
	fogBrokerPort string
}

func NewEnv(ctx context.Context, t *testing.T) (*Env, error) {
	cloudBroker, err := dependencies.NewMosquitto(ctx)
	if err != nil {
		return &Env{}, err
	}
	fogBroker, err := dependencies.NewMosquitto(ctx)
	if err != nil {
		return &Env{}, err
	}

	return &Env{
		cloudBroker: cloudBroker,
		fogBroker: fogBroker,
		UserID: "user",
	}, nil
}

func (e *Env) Start(ctx context.Context, t *testing.T, applicationLogChan chan string) (error) {
	t.Log("Start Verne")
	err, cloudBrokerPort := e.cloudBroker.StartAndWait(ctx)
	if err != nil {
		return err
	}
	e.cloudBrokerPort = cloudBrokerPort
	t.Log("Started Verne")

	t.Log("Start Mosquitto")
	err, fogBrokerPort := e.fogBroker.StartAndWait(ctx)
	if err != nil {
		return err
	}
	e.fogBrokerPort = fogBrokerPort
	t.Log("Started Mosquitto")

	t.Log("Start Connector")
	err = e.StartAndWait(ctx, t, applicationLogChan)
	if err != nil {
		return err
	}
	t.Log("Started Connector")
	
	return nil
}

func (e *Env) PublishToCloud(topic string, payload []byte, t *testing.T) error {
	e.platformMqttClient = mqtt.NewMQTTClient("localhost", e.cloudBrokerPort, nil, nil, false)
	t.Log("Connect to cloud broker at: " + e.cloudBrokerPort)
	err := e.platformMqttClient.ConnectMQTTBroker(nil, nil)
	if err != nil {
		return err
	}
	t.Log("Publish to: " + topic)
	return e.platformMqttClient.Publish(topic, string(payload), 2)
}

func (e *Env) PublishToFog(topic string, payload []byte, t *testing.T) error {
	e.fogMqttClient = mqtt.NewMQTTClient("localhost", e.fogBrokerPort, nil, nil, false)
	t.Log("Connect to local broker at " + e.fogBrokerPort)
	err := e.fogMqttClient.ConnectMQTTBroker(nil, nil)
	if err != nil {
		return err
	}
	t.Log("Publish to: " + topic)
	return e.fogMqttClient.Publish(topic, string(payload), 2)

}