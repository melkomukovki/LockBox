// Package grpcclient - описывает GRPC клиент и взаимодействие с бэкендом приложения
package grpcclient

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/melkomukovki/LockBox/api/pb"
)

// Client - структура GRPC клиента
type Client struct {
	conn          *grpc.ClientConn
	UserService   pb.UserServiceClient
	SecretService pb.SecretServiceClient
}

// NewGRPCClient - конструктор для получения экземпляра GRPC клиента
// Обязательный параметр - адрес сервера в формате "{host}:{port}"
func NewGRPCClient(serverAddress string) (*Client, error) {
	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	return &Client{
		conn:          conn,
		UserService:   pb.NewUserServiceClient(conn),
		SecretService: pb.NewSecretServiceClient(conn),
	}, nil
}

func (c *Client) withAuthContext(ctx context.Context, token string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %s", token))
}

func (c *Client) SignIn(ctx context.Context, login, password string) (*pb.SignInResponse, error) {
	return c.UserService.SignIn(ctx, &pb.SignInRequest{Login: login, Password: password})
}

func (c *Client) SignUp(ctx context.Context, login, password string) (*pb.SignUpResponse, error) {
	return c.UserService.SignUp(ctx, &pb.SignUpRequest{Login: login, Password: password})
}

func (c *Client) SecretsList(ctx context.Context, token string) (*pb.ListResponse, error) {
	ctx = c.withAuthContext(ctx, token)
	req := &pb.ListRequest{}
	return c.SecretService.List(ctx, req)
}

func (c *Client) SecretsGet(ctx context.Context, token string, id int64) (*pb.GetResponse, error) {
	ctx = c.withAuthContext(ctx, token)
	req := &pb.GetRequest{Id: id}
	return c.SecretService.Get(ctx, req)
}

func (c *Client) SecretsDelete(ctx context.Context, token string, id int64) (*pb.DeleteResponse, error) {
	ctx = c.withAuthContext(ctx, token)
	req := &pb.DeleteRequest{Id: id}
	return c.SecretService.Delete(ctx, req)
}

func (c *Client) SecretsAdd(ctx context.Context, token string, secret *pb.Secret) (*pb.StoreResponse, error) {
	ctx = c.withAuthContext(ctx, token)
	req := &pb.StoreRequest{Secret: secret}
	return c.SecretService.Store(ctx, req)
}

func (c *Client) Close() error {
	return c.conn.Close()
}
