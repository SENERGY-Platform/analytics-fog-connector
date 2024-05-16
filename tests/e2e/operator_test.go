package e2e

import (
	"context"
	"testing"
	"time"
	"github.com/hahahannes/e2e-go-utils/lib/streaming/mqtt"
	"github.com/stretchr/testify/assert"
)

func TestForwardStart(t *testing.T) {
	// Forward the operator start command
	ctx := context.Background()
	env, err := NewEnv(ctx, t)
	if err != nil {
		t.Error(err)
		return 
	}

	err = env.Start(ctx, t, make(chan string))
	if err != nil {
		t.Error(err)
		return 
	}

	operatorTopic :=  "fog/analytics/" + env.UserID + "/operator/control/start"
	fogOperatorTopic := "analytics/operator/control/start"
	msg := "test"

	result, err := mqtt.WaitForMQTTMessageReceived(fogOperatorTopic, ".*" + msg + ".*", func(context.Context) error {
		return env.PublishToCloud(operatorTopic, []byte(msg), t)
	}, 15 * time.Second, "localhost", env.fogBrokerPort, true)
	if err != nil {
		t.Error(err)
		return 
	}
	if result.Error != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, result.Received, true)
}

func TestForwardStop(t *testing.T) {
	// Forward the operator stop command
	ctx := context.Background()
	env, err := NewEnv(ctx, t)
	if err != nil {
		t.Error(err)
		return 
	}

	err = env.Start(ctx, t, make(chan string))
	if err != nil {
		t.Error(err)
		return 
	}

	operatorTopic :=  "fog/analytics/" + env.UserID + "/operator/control/stop"
	fogOperatorTopic := "analytics/operator/control/stop"
	msg := "test"

	result, err := mqtt.WaitForMQTTMessageReceived(fogOperatorTopic, ".*" + msg + ".*", func(context.Context) error {
		return env.PublishToCloud(operatorTopic, []byte(msg), t)
	}, 15 * time.Second, "localhost", env.fogBrokerPort, true)
	if err != nil {
		t.Error(err)
		return 
	}
	if result.Error != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, result.Received, true)
}

func TestOperatorSync(t *testing.T) {
	// Forward the operator sync result 

	ctx := context.Background()
	env, err := NewEnv(ctx, t)
	if err != nil {
		t.Error(err)
		return 
	}

	err = env.Start(ctx, t, make(chan string))
	if err != nil {
		t.Error(err)
		return 
	}

	operatorTopic :=  "fog/analytics/" + env.UserID + "/operator/control/sync/response"
	fogOperatorTopic := "analytics/operator/control/sync/response"
	msg := "test"

	result, err := mqtt.WaitForMQTTMessageReceived(fogOperatorTopic, ".*" + msg + ".*", func(context.Context) error {
		return env.PublishToCloud(operatorTopic, []byte(msg), t)
	}, 15 * time.Second, "localhost", env.fogBrokerPort, true)
	if err != nil {
		t.Error(err)
		return 
	}
	if result.Error != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, result.Received, true)
}