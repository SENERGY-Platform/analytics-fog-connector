package e2e

import (
	"testing"
	"context"
	"time"
	"github.com/hahahannes/e2e-go-utils/lib/mqtt"
)

func TestForwardCloudOperatorMsg(t *testing.T) {
	ctx := context.TODO()
	env, err := NewEnv(ctx, t)
	if err != nil {
		t.Error(err)
		return 
	}

	err, _, mosquittoPort := env.Start(ctx, t)
	if err != nil {
		t.Error(err)
		return 
	}

	operatorTopic :=  "/fog/analytics/" + env.UserID + "/downstream/operator"
	msg := "test"
	mqtt.WaitForMQTTMessageReceived("", msg, func() {
		env.PublishToCloud(operatorTopic, []byte(msg))
	}, 30 * time.Second, "localhost", mosquittoPort)
}

