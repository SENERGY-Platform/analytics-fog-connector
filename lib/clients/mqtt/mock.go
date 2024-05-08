package mqtt

type MockMqttClient struct {
	ReturnError bool;
	Error error;
}

func (m MockMqttClient) Publish(topic string, message string, qos int) error {
	if m.ReturnError {
		return m.Error
	}
	return nil
}

func (m MockMqttClient) Subscribe(topic string, qos int) error {
	if m.ReturnError {
		return m.Error
	}
	return nil
}

func (m MockMqttClient) Unsubscribe(topic string) error {
	if m.ReturnError {
		return m.Error
	}
	return nil
}