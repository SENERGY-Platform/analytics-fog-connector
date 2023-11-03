/*
 * Copyright 2019 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/SENERGY-Platform/analytics-fog-connector/lib/config"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/connector"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/logging"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/mqtt"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/relay"
	srv_base "github.com/SENERGY-Platform/go-service-base/util"
	"github.com/SENERGY-Platform/go-service-base/watchdog"

	"github.com/joho/godotenv"
)

func main() {
	ec := 0
	defer func() {
		os.Exit(ec)
	}()

	log.Println("Load .env file")
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file")
	}

	log.Println("Load config")
	config, err := config.NewConfig("")
	if err != nil {
		fmt.Println(err)
	}

	log.Println("Init Logger")
	logFile, err := logging.InitLogger(config.Logger)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		var logFileError *srv_base.LogFileError
		if errors.As(err, &logFileError) {
			ec = 1
			return
		}
	}
	if logFile != nil {
		defer logFile.Close()
	}

	watchdog.Logger = logging.Logger
	watchdog := watchdog.New(syscall.SIGINT, syscall.SIGTERM)

	// TODO userID
	userID := "dd69ea0d-f553-4336-80f3-7f4567f85c7b"
	mqttClient := mqtt.NewMQTTClient(config.Broker, userID, logging.Logger)

	// TODO connector
	connector := connector.NewConnector(mqttClient)

	relayController := relay.NewRelayController(connector, userID)

	mqttClient.ConnectMQTTBroker(relayController)

	watchdog.RegisterStopFunc(func() error {
		mqttClient.CloseConnection()
		return nil
	})

	watchdog.Start()

	ec = watchdog.Join()
}
