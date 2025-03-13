package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/melkomukovki/LockBox/internal/models"
	"github.com/melkomukovki/LockBox/internal/server/config"
	"github.com/melkomukovki/LockBox/internal/server/grpcserver"
	"github.com/melkomukovki/LockBox/internal/server/respository/postgres"
	"github.com/melkomukovki/LockBox/internal/server/service"
	"github.com/melkomukovki/LockBox/pkg/auth"
)

const timeFormat = "02/Jan/2006 15:04:05 -0700"

var (
	buildVersion = "N/A"
	buildCommit  = "N/A"
	buildDate    = "N/A"
)

var (
	configPath = flag.String("config", "config/config.yaml", "path to config file")
)

func init() {
	zerolog.TimeFieldFormat = timeFormat
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
}

func main() {
	printInfo()

	flag.Parse()

	// Load config
	cfg, err := config.New(*configPath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config file")
	}
	setLogLevel(cfg.Log.Level)

	log.Info().Msg("config loaded")
	log.Debug().Msgf("config: %+v", cfg)

	// Create DB connection & migrate
	pgDb, err := postgres.NewPgClient(&cfg.Postgres)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to postgres")
	}
	if err := pgDb.AutoMigrate(&models.User{}, &models.Secret{}); err != nil {
		log.Fatal().Err(err).Msg("failed to migrate models")
	}

	// Create repository
	userRepository := postgres.NewUserRepository(pgDb)
	secretRepository := postgres.NewSecretRepository(pgDb)

	// Create service
	userService := service.NewUserService(userRepository)
	secretService := service.NewSecretService(secretRepository)

	// Create JWT manager
	jwtManager, err := auth.NewManager(cfg.Auth.SigningKey, cfg.Auth.AccessTokenTTL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize jwt manager")
	}

	// Create grpc server
	srv := grpcserver.New(&cfg.Server, userService, secretService, jwtManager)

	// Run server
	go func() {
		if err := srv.Run(); err != nil {
			log.Fatal().Err(err).Msg("error while running grpc server")
		}
	}()
	log.Info().Msg("grpc server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		err := srv.Stop(ctx)
		if err != nil {
			log.Fatal().Err(err).Msg("error while stopping grpc server")
		}
	}()
}

func printInfo() {
	fmt.Printf("Version: %s\n", buildVersion)
	fmt.Printf("Commit: %s\n", buildCommit)
	fmt.Printf("Build date: %s\n", buildDate)
}

func setLogLevel(level string) {
	logLevel := map[string]zerolog.Level{
		"debug": zerolog.DebugLevel,
		"info":  zerolog.InfoLevel,
		"warn":  zerolog.WarnLevel,
		"error": zerolog.ErrorLevel,
		"fatal": zerolog.FatalLevel,
		"panic": zerolog.PanicLevel,
	}

	l, ok := logLevel[level]
	if ok {
		zerolog.SetGlobalLevel(l)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
