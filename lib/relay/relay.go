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
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/connector"

	"github.com/SENERGY-Platform/analytics-fog-lib/lib/control"
	upstreamLib "github.com/SENERGY-Platform/analytics-fog-lib/lib/upstream"
	downstreamLib "github.com/SENERGY-Platform/analytics-fog-lib/lib/downstream"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/logging"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type RelayController struct {
	Connector *connector.Connector
	UserID    string
	PublishResultsToPlatform bool
}

func NewRelayController(connector *connector.Connector, userID string, publishResultsToPlatform bool) *RelayController {
	return &RelayController{
		Connector: connector,
		UserID:    userID,
		PublishResultsToPlatform: publishResultsToPlatform,
	}
}

func (relay *RelayController) ProcessMessage(message MQTT.Message) {
	payload := message.Payload()
	switch message.Topic() {
	case control.GetConnectorControlTopic(relay.UserID):
		relay.processOperatorControlCommand(payload)
	case upstreamLib.UpstreamProxyDisableTopic:
		relay.processUpstreamDisable(payload)
	case upstreamLib.UpstreamProxyEnableTopic:
		relay.processUpstreamEnable(payload)
	case downstreamLib.GetDownstreamTopic(relay.UserID):
		relay.ForwardCloudMessageToFog(payload)
	default:
		// default are all operator topics that connector subscribed to
		relay.processMessageToUpstream(payload, message.Topic())
	}
}

func (relay *RelayController) OnMessageReceived(client MQTT.Client, message MQTT.Message) {
	logging.Logger.Debugf("Received message on topic: %s\nMessage: %s\n", message.Topic(), message.Payload())
	go relay.ProcessMessage(message)
}
