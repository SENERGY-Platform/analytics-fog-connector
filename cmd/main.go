/*
 * Copyright 2019 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"log"
	"os"

	"github.com/SENERGY-Platform/analytics-fog-connector/lib"
	"github.com/joho/godotenv"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/clients/auth"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/config"
)

func main() {
	ec := 0
	defer func() {
		os.Exit(ec)
	}()

	log.Println("Load .env file")
	_ = godotenv.Load()

	log.Println("Load config")
	config, err := config.NewConfig("")
	if err != nil {
		log.Println(err)
		ec = 1
	}
	authClient := auth.NewAuthClient(config.KeyCloakURL, config.ClientID)

	err = lib.Run(context.Background(), os.Stdout, os.Stderr, authClient, *config)
	if err != nil {
		log.Print("Error starting app: " + err.Error())
		ec = 1
	}
}
