package handlers

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/melkomukovki/LockBox/api/pb"
	"github.com/melkomukovki/LockBox/internal/models"
)

// SecretController описание структуры контроллера
type SecretController struct {
	pb.UnimplementedSecretServiceServer
	service models.SecretService
}

// NewSecretController конструктор для получения Secret контроллера
func NewSecretController(service models.SecretService) *SecretController {
	return &SecretController{service: service}
}

// Store - функция обработчик запросов на сохранение секретов
func (s *SecretController) Store(ctx context.Context, request *pb.StoreRequest) (*pb.StoreResponse, error) {

	secret := &models.Secret{
		Name:        request.Secret.Name,
		UserID:      ctx.Value("userId").(int),
		Description: request.Secret.Description,
		Type:        models.SecretType(request.Secret.Type),
		Data:        request.Secret.Data,
	}

	secretId, err := s.service.CreateSecret(ctx, secret)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.StoreResponse{Message: "secret created", Id: int64(secretId)}, nil
}

// List - функция обработчик для получения списка секретов пользователя
func (s *SecretController) List(ctx context.Context, request *pb.ListRequest) (*pb.ListResponse, error) {
	userId := ctx.Value("userId").(int)
	secrets, err := s.service.GetAllSecrets(ctx, userId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var pbSecrets []*pb.Secret
	for _, secret := range secrets {
		pbSecret := &pb.Secret{
			Id:          int64(secret.ID),
			Name:        secret.Name,
			Description: secret.Description,
			Type:        string(secret.Type),
			Data:        secret.Data,
		}
		pbSecrets = append(pbSecrets, pbSecret)
	}
	return &pb.ListResponse{Secrets: pbSecrets}, nil
}

// Delete - функция обработчик для запросов удаления секрета
func (s *SecretController) Delete(ctx context.Context, request *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	userId := ctx.Value("userId").(int)

	if err := s.service.DeleteSecret(ctx, int(request.Id), userId); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.DeleteResponse{Message: "secret deleted"}, nil
}

// Update - функция обработчик для запросов обновления секрета
func (s *SecretController) Update(ctx context.Context, request *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	secret := &models.Secret{
		ID:          int(request.Secret.Id),
		Name:        request.Secret.Name,
		Description: request.Secret.Description,
		Type:        models.SecretType(request.Secret.Type),
		Data:        request.Secret.Data,
	}

	userId := ctx.Value("userId").(int)
	if err := s.service.UpdateSecret(ctx, secret, userId); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.UpdateResponse{Message: "secret updated"}, nil
}

// Get - функция обработчик для получения информация по секрету
func (s *SecretController) Get(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	userId := ctx.Value("userId").(int)
	secret, err := s.service.GetSecret(ctx, int(request.Id), userId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbSecret := &pb.Secret{
		Id:          int64(secret.ID),
		Name:        secret.Name,
		Description: secret.Description,
		Type:        string(secret.Type),
		Data:        secret.Data,
	}

	return &pb.GetResponse{Secret: pbSecret}, nil
}
