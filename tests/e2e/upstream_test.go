package e2e

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/SENERGY-Platform/analytics-fog-lib/lib/upstream"
	testLib "github.com/hahahannes/e2e-go-utils/lib"
	"github.com/hahahannes/e2e-go-utils/lib/streaming/mqtt"
	"github.com/stretchr/testify/assert"
)

func TestDisableForwarding(t *testing.T) {
}

func TestSyncResponseForwarding(t *testing.T) {
	// Sync the forwarding rules based on a sync response from cloud
	// The sync is requested at startup and reconnects but here I just send the sync response manually
	
	ctx := context.Background()
	env, err := NewEnv(ctx, t)
	if err != nil {
		t.Error(err)
		return 
	}

	appLogChan := make(chan string)
	err = env.Start(ctx, t, appLogChan)
	if err != nil {
		t.Error(err)
		return 
	}

	localOperatorOutput := "Test"
	opName := "op"
	opID := "op1"
	pipeID := "pipe"
	localOperatorOutputTopic := "foo/bar/" + opName + "/" + opID + "/" + pipeID
	cloudOperatorOutputTopic := "fog/analytics/upstream/messages/analytics-" + opName
	syncResponseTopic := "fog/analytics/"+ env.UserID+ "/upstream/sync/response"

	result, err := mqtt.WaitForMQTTMessageReceived(cloudOperatorOutputTopic, ".*" + localOperatorOutput + ".*", func(context.Context) error {
		synCmd := upstream.UpstreamSyncMessage{
			OperatorOutputTopics: []string{localOperatorOutputTopic},
		}
		syncMessage, err := json.Marshal(synCmd)
		if err != nil {
			return err
		}
		received, err := testLib.WaitForStringReceived(".*Successfully subscribed to:.*", func (sendCtx context.Context) error {
			return env.PublishToCloud(syncResponseTopic, []byte(syncMessage), t)
		}, appLogChan, 30 * time.Second, true)
		if err != nil {
			return err
		}
		if received.Error != nil {
			return err
		}
		if received.Received == false {
			return errors.New("Subscribe Log to operator topic not received")
		}
		return env.PublishToFog(localOperatorOutputTopic, []byte(localOperatorOutput), t)
	}, 60 * time.Second, "localhost", env.cloudBrokerPort, false)
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
	err = env.Start(ctx, t, appLogChan)
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