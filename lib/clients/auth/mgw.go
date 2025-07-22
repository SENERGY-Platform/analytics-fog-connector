package auth

import (
	"context"
	"github.com/SENERGY-Platform/mgw-cloud-proxy/cert-manager/lib/client"
	"sync"
	"time"
)

type MgwIdpClient struct {
	client                *client.Client
	lastUserId            string
	lastUserIdRequestTime time.Time
	mux                   sync.Mutex
}

func (this *MgwIdpClient) GetUserID() (string, error) {
	this.mux.Lock()
	defer this.mux.Unlock()
	if this.lastUserId != "" && time.Since(this.lastUserIdRequestTime) < 10*time.Minute {
		return this.lastUserId, nil
	}
	network, err := this.client.NetworkInfo(context.Background(), "")
	if err != nil {
		return "", err
	}
	this.lastUserId = network.UserID
	this.lastUserIdRequestTime = time.Now()
	return this.lastUserId, nil
}
