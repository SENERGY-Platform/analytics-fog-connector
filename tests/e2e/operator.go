package e2e

import (
	"testing"
	"context"
)

func TestForwardStart(t *testing.T) {
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

func TestForwardStop(t *testing.T) {
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

func TestSyncResponse(t *testing.T) {
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
