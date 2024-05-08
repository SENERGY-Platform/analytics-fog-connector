package mqtt

type MQTTClient interface {
	Publish(topic string, message string, qos int) error
	Subscribe(topic string, qos int) error
	Unsubscribe(topic string) error 
}