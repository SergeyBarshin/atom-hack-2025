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

	portServ := config.Portserv
	grpcPort := config.GrpcPort

	fmt.Println("Server running on port", portServ)
	http.ListenAndServe(portServ, nil)

	// Настройка gRPC сервера
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	authent.RegisterAuthServiceServer(grpcServer, &server.AuthServer{})

	fmt.Println("gRPC Auth Service is running on port", grpcPort)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
