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
	"strings"

	"github.com/SENERGY-Platform/analytics-fog-connector/lib/connector"

	"github.com/SENERGY-Platform/analytics-fog-connector/lib/logging"
	downstreamLib "github.com/SENERGY-Platform/analytics-fog-lib/lib/downstream"
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/operator"
	upstreamLib "github.com/SENERGY-Platform/analytics-fog-lib/lib/upstream"

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
	topic := message.Topic() 

	switch topic {
	// Operator Control 
	case operator.GetStartOperatorCloudTopic(relay.UserID):
		relay.processStartOperatorCommand(payload)
		return
	case operator.GetStopOperatorCloudTopic(relay.UserID):
		relay.processStopOperatorCommand(payload)
		return
	case operator.GetOperatorControlSyncResponseTopic(relay.UserID):
		relay.processOperatorSync(payload)
		return

	// Upstream Forward Control
	case upstreamLib.GetUpstreamDisableCloudTopic(relay.UserID):
		relay.processUpstreamDisable(payload)
		return 
	case upstreamLib.GetUpstreamEnableCloudTopic(relay.UserID):
		relay.processUpstreamEnable(payload)
		return 
	case upstreamLib.GetUpstreamControlSyncResponseTopic(relay.UserID):
		relay.processUpstreamSync(payload)
		return
	}

	if strings.HasPrefix(topic, downstreamLib.GetDownstreamOperatorCloudMatchTopic(relay.UserID)) {
		// Prefix match
		relay.processOperatorDownstreamMessage(payload, topic)
		return
	}

	// default are all operator topics that connector subscribed to
	relay.processMessageToUpstream(payload, message.Topic())
}

func (relay *RelayController) OnMessageReceived(client MQTT.Client, message MQTT.Message) {
	logging.Logger.Debug("Received message on topic: " + message.Topic() + "\nMessage: " + string(message.Payload()))
	go relay.ProcessMessage(message)
}
