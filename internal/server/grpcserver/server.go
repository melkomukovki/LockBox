package grpcserver

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/melkomukovki/LockBox/api/pb"
	"github.com/melkomukovki/LockBox/internal/models"
	"github.com/melkomukovki/LockBox/internal/server/config"
	"github.com/melkomukovki/LockBox/internal/server/grpcserver/handlers"
	"github.com/melkomukovki/LockBox/internal/server/grpcserver/interceptors"
	"github.com/melkomukovki/LockBox/pkg/auth"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"time"
)

type Server struct {
	grpcServer *grpc.Server
	config     *config.ServerConfig
}

func New(cfg *config.ServerConfig, userService models.IUserService, secretService models.ISecretService, jwtManager auth.JWTManager) *Server {

	authInterceptor := interceptors.NewAuthInterceptor(jwtManager)

	server := grpc.NewServer(
		grpc.ConnectionTimeout(10*time.Second),
		grpc.MaxRecvMsgSize(1024*1024*16),
		grpc.MaxSendMsgSize(1024*1024*16),
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(interceptors.InterceptorLogger(log.Logger)),
			authInterceptor.UnaryInterceptor(),
		),
	)
	reflection.Register(server)

	userController := handlers.NewUserController(userService, jwtManager)
	secretController := handlers.NewSecretController(secretService)

	pb.RegisterUserServiceServer(server, userController)
	pb.RegisterSecretServiceServer(server, secretController)

	return &Server{grpcServer: server, config: cfg}
}

func (s *Server) Run() error {
	address := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	return s.grpcServer.Serve(lis)
}

func (s *Server) Stop(ctx context.Context) error {
	ok := make(chan struct{})
	go func() {
		s.grpcServer.GracefulStop()
		close(ok)
	}()

	select {
	case <-ok:
		return nil
	case <-ctx.Done():
		s.grpcServer.Stop()
		return ctx.Err()
	}
}
