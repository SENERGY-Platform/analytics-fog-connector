package e2e

import (
	"context"
	"testing"
	"time"
	"github.com/hahahannes/e2e-go-utils/lib/streaming/mqtt"
	"github.com/stretchr/testify/assert"
)

func TestDownstream(t *testing.T) {
	ctx := context.Background()
	env, err := NewEnv(ctx, t)
	if err != nil {
		t.Error(err)
		return 
	}

	err = env.Start(ctx, t, make(chan string), "downstream")
	if err != nil {
		t.Error(err)
		return 
	}
	t.Log("Run test")

	operatorName := "operator"
	operatorID := "op1"
	pipelineID := "pipeline"
	operatorTopic :=  "fog/analytics/" + env.UserID + "/downstream/operator/" + operatorName + "/" + operatorID + "/" + pipelineID
	fogOperatorTopic := "analytics/" + operatorName + "/" + operatorID + "/" + pipelineID
	msg := "test"

	mqttCtx, mqttCf := context.WithTimeout(ctx, 30 * time.Second)
	defer mqttCf()
	
	result, err := mqtt.WaitForMQTTMessageReceived(mqttCtx, fogOperatorTopic, ".*" + msg + ".*", func(context.Context) error {
		return env.PublishToCloud(operatorTopic, []byte(msg), t)
	}, "localhost", env.fogBrokerPort, false)
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

