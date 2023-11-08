module github.com/SENERGY-Platform/analytics-fog-connector

go 1.21.3

replace github.com/SENERGY-Platform/analytics-fog-lib => ../analytics-fog-lib

require (
	github.com/Nerzal/gocloak/v13 v13.8.0
	github.com/SENERGY-Platform/analytics-fog-lib v1.0.2
	github.com/SENERGY-Platform/go-service-base/util v0.14.0
	github.com/SENERGY-Platform/go-service-base/watchdog v0.4.1
	github.com/eclipse/paho.mqtt.golang v1.4.3
	github.com/joho/godotenv v1.5.1
	github.com/y-du/go-log-level v0.2.3
)

require (
	github.com/go-resty/resty/v2 v2.7.0 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/segmentio/ksuid v1.0.4 // indirect
	github.com/y-du/go-env-loader v0.5.1 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sync v0.4.0 // indirect
)
