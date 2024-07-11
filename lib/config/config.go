package config

import (
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/mqtt"
	srv_base "github.com/SENERGY-Platform/go-service-base/util"
)

type Config struct {
	Debug bool `json:"debug" env_var:"DEBUG"`
	FogBroker      mqtt.FogBrokerConfig
	PlatformBroker mqtt.PlatformBrokerConfig
	DataDir     string                `json:"data_dir" env_var:"DATA_DIR"`
	KeyCloakURL string                `json:"keycloak_url" env_var:"KEYCLOAK_URL"`
	ClientID    string                `json:"client_id" env_var:"CLIENT_ID"`
	Username    string                `json:"username" env_var:"AUTH_USERNAME"`
	Password    string                `json:"password" env_var:"AUTH_PASSWORD"`
	PublishResultsToPlatform bool 	`json:"publish_to_platform" env_var:"PUBLISH_RESULTS_TO_PLATFORM"`
	SyncIntervalInSeconds    int                `json:"sync_interval" env_var:"SYNC_INTERVAL"`
}

func NewConfig(path string) (*Config, error) {
	cfg := Config{
		FogBroker: mqtt.FogBrokerConfig{
			Port: "1883",
			Host: "localhost",
		},
		Debug: false,
		DataDir: "./data",
	}

	err := srv_base.LoadConfig(path, &cfg, nil, nil, nil)
	return &cfg, err
}
