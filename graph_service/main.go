package main

import (
	"database/sql"
	"graph_service/handlers"
	"graph_service/repositories"
	"graph_service/service"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {

	port := os.Getenv("PORT")
	dbUrl := os.Getenv("DB_URL")

	//Подключение к базе данных
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err) //вообще работать не будет
	}
	defer db.Close()

	//Создание экземпляра для работы с бд
	graphRepo := repositories.NewGraphRepositoryPostgres(db)

	//Создание экземпляра сервиса данных
	dataService := service.NewDataService(graphRepo)

	//Создание экземпляров обработчиков HTTP-запросов
	graphHandler := handlers.NewGraphHandler(dataService)

	router := gin.Default()

	router.POST("/create_graph", graphHandler.CreateGraph)
	router.GET("/get_graph", graphHandler.GetGraph)
	router.PUT("/edit_graph", graphHandler.UpdateGraph)
	router.DELETE("/delete_graphs", graphHandler.DeleteGraph)

	router.Run(":" + port)
}
