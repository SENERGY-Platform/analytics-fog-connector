package e2e

import (
	"context"
	"testing"
	"time"
	"github.com/hahahannes/e2e-go-utils/lib/streaming/mqtt"
	"github.com/stretchr/testify/assert"
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/upstream"
	"encoding/json"
	testLib "github.com/hahahannes/e2e-go-utils/lib"
)

func TestDisableForwarding(t *testing.T) {

}

func TestSyncResponseForwarding(t *testing.T) {
	// Sync the enabled/disabled forwarding for local operator output 
	
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

	operatorTopic :=  "fog/analytics/" + env.UserID + "/upstream/enable"
	fogOperatorTopic := "analytics/operator/control/stop"
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

func TestEnableForwarding(t *testing.T) {
	// Enabling of upstream forwarding and forwarding of local operator 

	ctx := context.Background()
	env, err := NewEnv(ctx, t)
	if err != nil {
		t.Error(err)
		return 
	}

	appLogChan := make(chan string)
	err, _, _ = env.Start(ctx, t, appLogChan)
	if err != nil {
		t.Error(err)
		return 
	}

	localOperatorOutput := "Test"
	opName := "op"
	opID := "op1"
	pipeID := "pipe"
	localOperatortopic := "foo/bar/" + opName + "/" + opID + "/" + pipeID
	cloudTopic := "fog/analytics/upstream/messages/analytics-" + opName
	result, err := mqtt.WaitForMQTTMessageReceived(cloudTopic, ".*" + localOperatorOutput + ".*", func(context.Context) error {
		EnableForwarding(t, env, localOperatortopic, appLogChan)
		return env.PublishToFog(localOperatortopic, []byte(localOperatorOutput), t)
	}, 15 * time.Second, "localhost", env.cloudBrokerPort, true)
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

func EnableForwarding(t *testing.T, env *Env, localOperatortopic string, appLogChan chan string) error {
	// Util function to wait for connector to successfully subscribe to operator topic
	received, err := testLib.WaitForStringReceived(".*Successfully subscribed to:.*", func (sendCtx context.Context) error {
		operatorTopic :=  "fog/analytics/" + env.UserID + "/upstream/enable"
		cmd := upstream.UpstreamControlMessage{
			OperatorOutputTopic: localOperatortopic,
		}
		msg, err := json.Marshal(cmd)
		if err != nil {
			return err
		}
		return env.PublishToCloud(operatorTopic, msg, t)
	}, appLogChan, 30 * time.Second, true)
	if err != nil {
		return err
	}
	if received.Error != nil {
		return err
	}
	return nil
}