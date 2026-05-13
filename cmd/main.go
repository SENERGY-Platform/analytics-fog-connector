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

	"github.com/SENERGY-Platform/analytics-fog-connector/internal/config"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib"
	"github.com/SENERGY-Platform/analytics-fog-connector/lib/clients/auth"
	"github.com/joho/godotenv"
)

func main() {
	ec := 0
	defer func() {
		os.Exit(ec)
	}()

	log.Println("Load .env file")
	_ = godotenv.Load()

	config.ParseFlags()

	log.Println("Load config")
	cfg, err := config.New(config.ConfPath)
	if err != nil {
		log.Println(err)
		ec = 1
		return
	}
	authClient := auth.NewAuthClient(cfg)

	err = lib.Run(context.Background(), os.Stdout, os.Stderr, authClient, *cfg)
	if err != nil {
		log.Print("Error starting app: " + err.Error())
		ec = 1
	}
}
