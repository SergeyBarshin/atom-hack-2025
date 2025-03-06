package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	SecretKey string
	DBURL     string
)

func LoadConfig() {
	godotenv.Load()
	SecretKey = os.Getenv("JWT_SECRET")
	DBURL = os.Getenv("DATABASE_URL")
	if SecretKey == "" || DBURL == "" {
		panic("Config values missing")
	}
	fmt.Println("Config loaded successfully")
}
