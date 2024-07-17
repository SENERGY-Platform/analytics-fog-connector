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

package subscriptionhandler

import (
	"strings"

	"github.com/SENERGY-Platform/analytics-fog-connector/lib/connector"

	"github.com/SENERGY-Platform/analytics-fog-connector/lib/logging"
	downstreamLib "github.com/SENERGY-Platform/analytics-fog-lib/lib/downstream"
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/operator"
	upstreamLib "github.com/SENERGY-Platform/analytics-fog-lib/lib/upstream"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type SubscriptionHandler struct {
	Connector *connector.Connector
	UserID    string
	PublishResultsToPlatform bool
}

func NewSubscriptionHandler(connector *connector.Connector, userID string, publishResultsToPlatform bool) *SubscriptionHandler {
	return &SubscriptionHandler{
		Connector: connector,
		UserID:    userID,
		PublishResultsToPlatform: publishResultsToPlatform,
	}
}

func (handler *SubscriptionHandler) ProcessMessage(message MQTT.Message) {
	payload := message.Payload()
	topic := message.Topic() 

	switch topic {
	// Operator Control 
	case operator.GetStartOperatorCloudTopic(handler.UserID):
		handler.processStartOperatorCommand(payload)
		return
	case operator.GetStopOperatorCloudTopic(handler.UserID):
		handler.processStopOperatorCommand(payload)
		return
	case operator.GetOperatorControlSyncResponseTopic(handler.UserID):
		handler.processOperatorSync(payload)
		return

	// Upstream Forward Control
	case upstreamLib.GetUpstreamDisableCloudTopic(handler.UserID):
		handler.processUpstreamDisable(payload)
		return 
	case upstreamLib.GetUpstreamEnableCloudTopic(handler.UserID):
		handler.processUpstreamEnable(payload)
		return 
	case upstreamLib.GetUpstreamControlSyncResponseTopic(handler.UserID):
		handler.processUpstreamSync(payload)
		return
	}

	if strings.HasPrefix(topic, downstreamLib.GetDownstreamOperatorCloudMatchTopic(handler.UserID)) {
		// Prefix match
		handler.processOperatorDownstreamMessage(payload, topic)
		return
	}

	// default are all operator topics that connector subscribed to
	handler.processMessageToUpstream(payload, message.Topic())
}

func (handler *SubscriptionHandler) OnMessageReceived(client MQTT.Client, message MQTT.Message) {
	logging.Logger.Debug("Received message on topic: " + message.Topic() + "\nMessage: " + string(message.Payload()))
	go handler.ProcessMessage(message)
}
