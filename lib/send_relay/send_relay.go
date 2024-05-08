package send_relay 

import (
	"errors"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/logging"
	"sync"
)

const logPrefix = "[send-relay]"

type SendRelayHandler struct {
	messages   chan lib.Message
	sendFunc   func(topic string, data []byte) error
	dChan      chan struct{}
	mu         sync.RWMutex
}

func New(buffer int, sendFunc func(topic string, data []byte) error) *SendRelayHandler {
	return &SendRelayHandler{
		messages:   make(chan lib.Message, buffer),
		sendFunc:   sendFunc,
		dChan:      make(chan struct{}),
	}
}

func (h *SendRelayHandler) Put(m lib.Message) error {
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
		topic := message.Payload()
		if err := h.sendFunc(message.Topic(), topic); err != nil {
			logging.Logger.Errorf("%s publish on topic (%s): %s", logPrefix, topic, err)
		}
	}
	h.dChan <- struct{}{}
}