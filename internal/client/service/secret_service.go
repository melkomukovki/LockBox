package service

import (
	"context"
	"errors"

	"github.com/melkomukovki/LockBox/api/pb"
	"github.com/melkomukovki/LockBox/internal/client/grpcclient"
)

// SecretService структура Secret сервиса
type SecretService struct {
	conn *grpcclient.Client
}

// List функция для получения списка секретов пользователя
func (s *SecretService) List(ctx context.Context, token string) ([]*pb.Secret, error) {
	resp, err := s.conn.SecretsList(ctx, token)
	if err != nil {
		return nil, err
	}

	if len(resp.Secrets) == 0 {
		return []*pb.Secret{}, errors.New("no secrets found")
	}

	return resp.Secrets, nil
}

// Get функция для получения конкретного секрета
func (s *SecretService) Get(ctx context.Context, token string, id int64) (*pb.Secret, error) {

	resp, err := s.conn.SecretsGet(ctx, token, id)
	if err != nil {
		return nil, err
	}

	return resp.Secret, nil

}

// Delete функция для удаления конкретного секрета
func (s *SecretService) Delete(ctx context.Context, token string, id int64) error {
	if _, err := s.conn.SecretsDelete(ctx, token, id); err != nil {
		return err
	}

	return nil
}

// Add функция для добавления секрета
func (s *SecretService) Add(ctx context.Context, token string, secret *pb.Secret) (int64, error) {
	resp, err := s.conn.SecretsAdd(ctx, token, secret)
	if err != nil {
		return 0, err
	}
	return resp.Id, nil
}

// NewSecretService - функция конструктор для получения экземпляра Secret сервиса
func NewSecretService(conn *grpcclient.Client) (*SecretService, error) {
	if conn == nil {
		return nil, errors.New("GRPC client must not be nil")
	}
	return &SecretService{conn: conn}, nil
}
