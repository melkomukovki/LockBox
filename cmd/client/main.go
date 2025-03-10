package main

import (
	"context"
	"fmt"
	"github.com/melkomukovki/LockBox/internal/client/app"
	"github.com/melkomukovki/LockBox/internal/client/grpcclient"
	"github.com/melkomukovki/LockBox/internal/client/service"
	"log"
	"os"
	"os/signal"
	"syscall"
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

	grpc, err := grpcclient.NewGRPCClient(serverAddress)
	if err != nil {
		log.Fatalf("cannot create grpc client: %v", err)
	}

	uService := service.NewUserService(grpc)
	sService := service.NewSecretService(grpc)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	client := app.NewApp(uService, sService)
	_ = client.Run(ctx, os.Args)
}

func printInfo() {
	fmt.Printf("Version: %s\n", buildVersion)
	fmt.Printf("Commit: %s\n", buildCommit)
	fmt.Printf("Build date: %s\n", buildDate)
}
