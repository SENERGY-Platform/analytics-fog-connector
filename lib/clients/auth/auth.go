package auth

import (
	"net/http"

	"github.com/Nerzal/gocloak/v13"
	"github.com/SENERGY-Platform/analytics-fog-connector/internal/config"
	"github.com/SENERGY-Platform/mgw-cloud-proxy/cert-manager/lib/client"
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
