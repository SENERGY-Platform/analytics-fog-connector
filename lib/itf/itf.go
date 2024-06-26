package itf

import "time"

type Message interface {
	Topic() string
	Payload() []byte
	Timestamp() time.Time
}

type MessageRelayHandler interface {
	Put(m Message) error
}

type MqttClient interface {
	Subscribe(topic string, qos byte, messageHandler func(m Message)) error
	Unsubscribe(topic string) error
	Publish(topic string, qos byte, retained bool, payload any) error
}