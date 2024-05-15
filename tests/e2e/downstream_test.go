package e2e

import (
	"context"
	"testing"
	"time"
	"github.com/hahahannes/e2e-go-utils/lib/streaming/mqtt"
	"github.com/stretchr/testify/assert"
)

func TestForwardCloudOperatorMsg(t *testing.T) {
	// Forward a message e.g. cloud operator output 
	ctx := context.Background()
	env, err := NewEnv(ctx, t)
	if err != nil {
		t.Error(err)
		return 
	}

	err, _, mosquittoPort := env.Start(ctx, t, make(chan string))
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

	result, err := mqtt.WaitForMQTTMessageReceived(fogOperatorTopic, ".*" + msg + ".*", func(context.Context) error {
		return env.PublishToCloud(operatorTopic, []byte(msg), t)
	}, 15 * time.Second, "localhost", mosquittoPort, true)
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

