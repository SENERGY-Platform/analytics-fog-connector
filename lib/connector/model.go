package connector

import "time"

type Message struct {
	topic string;
	payload []byte;
	timestamp time.Time
}

func(message Message) Topic() string {
	return message.topic
}

func(message Message) Payload() []byte {
	return message.Payload()
}

func(message Message) Timestamp() time.Time {
	return message.timestamp
}