/*
 * Copyright 2026 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package config

import (
	"github.com/SENERGY-Platform/analytics-fog-lib/lib/mqtt"
	sbconfighdl "github.com/SENERGY-Platform/go-service-base/config-hdl"
)

type Config struct {
	Debug                    bool `json:"debug" env_var:"DEBUG"`
	FogBroker                mqtt.FogBrokerConfig
	PlatformBroker           mqtt.PlatformBrokerConfig
	DataDir                  string `json:"data_dir" env_var:"DATA_DIR"`
	MgwCertManagerUrl        string `json:"mgw_cert_manager_url" env_var:"MGW_CERT_MANAGER_URL"` //alternative to KeyCloakURL
	KeyCloakURL              string `json:"keycloak_url" env_var:"KEYCLOAK_URL"`                 //may be ignored if MgwCertManagerUrl is used
	ClientID                 string `json:"client_id" env_var:"CLIENT_ID"`
	Username                 string `json:"username" env_var:"AUTH_USERNAME"`
	Password                 string `json:"password" env_var:"AUTH_PASSWORD"`
	PublishResultsToPlatform bool   `json:"publish_to_platform" env_var:"PUBLISH_RESULTS_TO_PLATFORM"`
	SyncIntervalInSeconds    int    `json:"sync_interval" env_var:"SYNC_INTERVAL"`
}

func New(path string) (*Config, error) {
	cfg := Config{
		FogBroker: mqtt.FogBrokerConfig{
			Port: "1883",
			Host: "localhost",
		},
		Debug:                 false,
		DataDir:               "./data",
		SyncIntervalInSeconds: 120,
	}

	err := sbconfighdl.Load(&cfg, nil, nil, nil, path)
	return &cfg, err
}
