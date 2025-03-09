package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/SergeyBarshin/atom-hack-2025/backend/auth-service/config"
	"github.com/SergeyBarshin/atom-hack-2025/backend/auth-service/db"
	"github.com/SergeyBarshin/atom-hack-2025/backend/auth-service/routes"
	"github.com/SergeyBarshin/atom-hack-2025/backend/auth-service/server"
	authent "github.com/SergeyBarshin/atom-hack-2025/backend/graph_service/grpc"

	"google.golang.org/grpc"
)

func main() {
	config.LoadConfig()
	db.InitDB()
	routes.SetupRoutes()

	port := config.Port
	grpcPort := config.GrpcPort

	fmt.Println("Server running on port", port)
	http.ListenAndServe(port, nil)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	// Создаем gRPC сервер
	grpcServer := grpc.NewServer()

	// Регистрируем AuthServer в gRPC
	authent.RegisterAuthServiceServer(grpcServer, &server.AuthServer{})

	fmt.Println("gRPC Auth Service is running on port", grpcPort)

	// Запускаем сервер
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
