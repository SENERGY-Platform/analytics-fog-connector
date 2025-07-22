package auth

import (
	gocloak "github.com/Nerzal/gocloak/v13"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/config"
	"github.com/SENERGY-Platform/mgw-cloud-proxy/cert-manager/lib/client"
	"net/http"
)

func NewAuthClient(config *config.Config) IDP {
	if config.MgwCertManagerUrl != "" {
		return &MgwIdpClient{client: client.New(http.DefaultClient, config.MgwCertManagerUrl)}
	}
	return &KeycloakClient{
		Client:   gocloak.NewClient(config.KeyCloakURL),
		ClientID: config.ClientID,
		config:   config,
	}
}
