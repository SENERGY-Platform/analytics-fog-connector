package e2e

import (
	"context"
	"testing"
)

func TestEnableForwarding(t *testing.T) {
	ctx := context.TODO()
	env, err := NewEnv(ctx, t)
	if err != nil {
		t.Error(err)
		return 
	}

	err = env.Start(ctx, t)
	if err != nil {
		t.Error(err)
		return 
	}
}

func TestDisableForwarding(t *testing.T) {
	ctx := context.TODO()
	env, err := NewEnv(ctx, t)
	if err != nil {
		t.Error(err)
		return 
	}

	err = env.Start(ctx, t)
	if err != nil {
		t.Error(err)
		return 
	}
}

func TestSyncResponseForwarding(t *testing.T) {
	ctx := context.TODO()
	env, err := NewEnv(ctx, t)
	if err != nil {
		t.Error(err)
		return 
	}

	err = env.Start(ctx, t)
	if err != nil {
		t.Error(err)
		return 
	}
}