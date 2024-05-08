package auth

import (
	"context"

	gocloak "github.com/Nerzal/gocloak/v13"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/logging"

)

type KeycloakAuthClient struct {
	Client   *gocloak.GoCloak
	ClientID string
}

func NewAuthClient(keycloakURL string, clientID string) *KeycloakAuthClient {
	return &KeycloakAuthClient{
		Client:   gocloak.NewClient(keycloakURL),
		ClientID: clientID,
	}
}

func (client *KeycloakAuthClient) GetUserID(username string, password string) (string, error) {
	ctx := context.Background()
	token, err := client.Client.Login(ctx, client.ClientID, "", "master", username, password)
	if err != nil {
		logging.Logger.Errorf("Cant login user %s: %s", username, err)
		return "", err
	}
	userInfo, err := client.Client.GetUserInfo(ctx, token.AccessToken, "master")
	if err != nil {
		logging.Logger.Errorf("Cant get user info from token: %s", err)
		return "", err
	}
	return *userInfo.Sub, nil
}
