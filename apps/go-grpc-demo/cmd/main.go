package main

import (
	"log"
	"os"

	grpcserver "dizzycoder1112/Dockerize-Monorepo-Structure-In-Node-And-Golang/grpc"
	hellopb "dizzycoder1112/Dockerize-Monorepo-Structure-In-Node-And-Golang/grpc/pb/hello"
	logger "dizzycoder1112/Dockerize-Monorepo-Structure-In-Node-And-Golang/logger"

	"go-grpc-demo/internal/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "50052"
	}

	appLogger := logger.NewConsole(logger.ConsoleOptions{
		ServiceName: "go-grpc-demo",
		Colored:     true,
	})

	server := grpcserver.NewServer(&grpcserver.ServerDeps{Log: appLogger})
	hellopb.RegisterGreeterServer(server.GrpcServer(), &handlers.Greeter{})
	server.EnableReflection()

	if err := server.Run(port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
	server.GracefulStop()
}
