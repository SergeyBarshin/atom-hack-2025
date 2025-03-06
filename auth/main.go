package main

import (
	"auth-service/config" // Подключаем наш конфиг
	"fmt"
)

func main() {
	config.LoadConfig() // Загружаем переменные окружения

	fmt.Println("JWT Secret:", config.SecretKey)
	fmt.Println("Database URL:", config.DBURL)
}
