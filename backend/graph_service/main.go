package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/SergeyBarshin/atom-hack-2025/backend/graph_service/handlers"
	"github.com/SergeyBarshin/atom-hack-2025/backend/graph_service/repositories"
	"github.com/SergeyBarshin/atom-hack-2025/backend/graph_service/service"

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

	router.POST("/graph/create_graph", graphHandler.CreateGraph)
	router.GET("/graph/get_graph", graphHandler.GetGraph)
	router.PUT("/graph/edit_graph", graphHandler.UpdateGraph)
	router.DELETE("/graph/delete_graph", graphHandler.DeleteGraph)
	router.GET("/graph/list_graphs", graphHandler.DeleteGraph)

	router.Run(":" + port)
}
