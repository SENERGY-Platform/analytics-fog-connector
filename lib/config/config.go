package config

import (
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/mqtt"
	srv_base "github.com/SENERGY-Platform/go-service-base/util"
	"github.com/y-du/go-log-level/level"
)

type Config struct {
	FogBroker      mqtt.FogBrokerConfig
	PlatformBroker mqtt.PlatformBrokerConfig
	Logger      srv_base.LoggerConfig `json:"logger" env_var:"LOGGER_CONFIG"`
	DataDir     string                `json:"data_dir" env_var:"DATA_DIR"`
	KeyCloakURL string                `json:"keycloak_url" env_var:"KEYCLOAK_URL"`
	ClientID    string                `json:"client_id" env_var:"CLIENT_ID"`
	Username    string                `json:"username" env_var:"AUTH_USERNAME"`
	Password    string                `json:"password" env_var:"AUTH_PASSWORD"`
}

func NewConfig(path string) (*Config, error) {
	cfg := Config{
		FogBroker: mqtt.FogBrokerConfig{
			Port: "1883",
			Host: "localhost",
		},
		Logger: srv_base.LoggerConfig{
			Level:        level.Debug,
			Utc:          true,
			Microseconds: true,
			Terminal:     true,
		},
		DataDir: "./data",
	}

	err := srv_base.LoadConfig(path, &cfg, nil, nil, nil)
	return &cfg, err
}
