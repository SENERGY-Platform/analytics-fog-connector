package auth

import (
	"context"
	"fmt"
	"github.com/Nerzal/gocloak/v13"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/config"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/logging"
)

type KeycloakClient struct {
	Client   *gocloak.GoCloak
	ClientID string
	config   *config.Config
}

func (client *KeycloakClient) GetUserID() (string, error) {
	return client.getUserID(client.config.Username, client.config.Password)
}

func (client *KeycloakClient) getUserID(username string, password string) (string, error) {
	ctx := context.Background()
	token, err := client.Client.Login(ctx, client.ClientID, "", "master", username, password)
	if err != nil {
		logging.Logger.Error(fmt.Sprintf("Cant login user %s: %s", username, err))
		return "", err
	}
	userInfo, err := client.Client.GetUserInfo(ctx, token.AccessToken, "master")
	if err != nil {
		logging.Logger.Error("Cant get user info from token: " + err.Error())
		return "", err
	}
	return *userInfo.Sub, nil
}
