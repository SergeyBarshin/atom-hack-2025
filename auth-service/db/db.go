package db

import (
	"database/sql"
	"log"

	"auth-service/config"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("postgres", config.DBURL)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Database not responding:", err)
	}

	log.Println("Connected to database")
}
