package send_relay

import (
	"errors"
	"sync"
	"fmt"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/itf"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/logging"
)

const logPrefix = "[send-relay]"

type SendRelayHandler struct {
	messages   chan itf.Message
	sendFunc   func(topic string, data []byte) error
	dChan      chan struct{}
	mu         sync.RWMutex
}

func New(buffer int, sendFunc func(topic string, data []byte) error) *SendRelayHandler {
	return &SendRelayHandler{
		messages:   make(chan itf.Message, buffer),
		sendFunc:   sendFunc,
		dChan:      make(chan struct{}),
	}
}

func (h *SendRelayHandler) Put(m itf.Message) error {
	select {
	case h.messages <- m:
	default:
		return errors.New("buffer full")
	}
	return nil
}

func (h *SendRelayHandler) Start() {
	go h.run()
}

func (h *SendRelayHandler) Stop() {
	close(h.messages)
	<-h.dChan
}

func (h *SendRelayHandler) run() {
	for message := range h.messages {
		topic := message.Topic()
		payload := message.Payload()
		logging.Logger.Debug("Try to send in order to topic: " + topic + " - Msg: " + string(payload))
		if err := h.sendFunc(topic, payload); err != nil {
			logging.Logger.Error(fmt.Sprintf("%s publish on topic (%s): %s", logPrefix, topic, err))
		}
	}
	h.dChan <- struct{}{}
}