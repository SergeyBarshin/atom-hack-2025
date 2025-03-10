package config

import (
	"fmt"
	"os"
)

var (
	SecretKey string
	DBURL     string
	Portlis   string
	GrpcPort  string
	Portserv  string
)

func LoadConfig() {
	SecretKey = os.Getenv("JWT_SECRET_AUTH")
	DBURL := os.Getenv("DB_URL_AUTH")
	Portserv = os.Getenv("PORT_AUTH")
	GrpcPort = os.Getenv("GRPC_PORT_AUTH")

	if SecretKey == "" || DBURL == "" || Portserv == "" || GrpcPort == "" {
		panic("Config values missing")
	}
	fmt.Println("Config loaded successfully")
}
