package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/melkomukovki/LockBox/internal/client/app"
	"github.com/melkomukovki/LockBox/internal/client/grpcclient"
	"github.com/melkomukovki/LockBox/internal/client/service"
	"github.com/melkomukovki/LockBox/internal/client/ui"
)

const (
	defaultServerAddress = "localhost:8888"
)

var (
	buildVersion = "N/A"
	buildCommit  = "N/A"
	buildDate    = "N/A"
)

func main() {
	printInfo()

	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress == "" {
		serverAddress = defaultServerAddress
	}

	// Создаем grpc клиент
	grpc, err := grpcclient.NewGRPCClient(serverAddress)
	if err != nil {
		log.Fatalf("cannot create grpc client: %v", err)
	}

	// Сервисный слой
	uService, err := service.NewUserService(grpc)
	if err != nil {
		log.Fatalf("cannot create user service: %v", err)
	}
	sService, err := service.NewSecretService(grpc)
	if err != nil {
		log.Fatalf("cannot create secret service: %v", err)
	}

	// Слой представления
	userHandler := ui.NewUserHandler(uService)
	secretHandler := ui.NewSecretHandler(sService)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	client := app.NewApp(userHandler, secretHandler)
	_ = client.Run(ctx, os.Args)
}

func printInfo() {
	fmt.Printf("Version: %s\n", buildVersion)
	fmt.Printf("Commit: %s\n", buildCommit)
	fmt.Printf("Build date: %s\n", buildDate)
}
