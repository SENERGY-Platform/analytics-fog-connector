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

package relay

import (
	"fmt"

	"github.com/SENERGY-Platform/analytics-fog-connector/lib/connector"
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/control"
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/operator"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type RelayController struct {
	Connector *connector.Connector
	UserID    string
}

func NewRelayController(connector *connector.Connector, userID string) *RelayController {
	return &RelayController{
		Connector: connector,
		UserID:    userID,
	}
}

func (relay *RelayController) ProcessMessage(message MQTT.Message) {
	switch message.Topic() {
	case control.GetConnectorControlTopic(relay.UserID):
		relay.processOperatorControlCommand(message.Payload())
	case operator.OperatorsResultTopic:
		relay.processOperatorOutputMessage(message.Payload())
	}

	// TODO Operator Output must be forwarded to platform
}

func (relay *RelayController) OnMessageReceived(client MQTT.Client, message MQTT.Message) {
	fmt.Printf("Received message on topic: %s\nMessage: %s\n", message.Topic(), message.Payload())
	go relay.ProcessMessage(message)
}
