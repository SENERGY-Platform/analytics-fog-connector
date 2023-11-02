package connector

import "github.com/SENERGY-Platform/analytics-fog-lib/lib/control"

func (controller *Connector) ForwardToMaster(message []byte) {
	controller.MQTTClient.Publish(control.ControlTopic, string(message), 2)
}
