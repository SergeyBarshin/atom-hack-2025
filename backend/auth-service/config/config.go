package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	SecretKey string
	DBURL     string
	Port      string
	GrpcPort  string
)

func LoadConfig() {
	godotenv.Load()
	SecretKey = os.Getenv("JWT_SECRET")
	DBURL := os.Getenv("DB_URL")
	Port = os.Getenv("PORT")
	GrpcPort = os.Getenv("GRPC_PORT")

	if SecretKey == "" || DBURL == "" || Port == "" || GrpcPort == "" {
		panic("Config values missing")
	}
	fmt.Println("Config loaded successfully")
}
