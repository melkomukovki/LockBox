package interceptors

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/melkomukovki/LockBox/pkg/auth"
)

// AuthInterceptor структура перехватчика авторизации пользователя
type AuthInterceptor struct {
	jwtManager auth.JWTManager
}

// NewAuthInterceptor конструктор для получения экземпляра перехватчика
func NewAuthInterceptor(jwtManager auth.JWTManager) *AuthInterceptor {
	return &AuthInterceptor{jwtManager: jwtManager}
}

// UnaryInterceptor проверяет содержимое метаданных в запросе на наличие JWT access токена
// Ожидаемый формат "authorization": "Bearer {access_token}"
func (a *AuthInterceptor) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if strings.Contains(info.FullMethod, "SecretService") {
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				return nil, status.Error(codes.Unauthenticated, "missing metadata")
			}

			authHeader := md["authorization"]
			if len(authHeader) == 0 {
				return nil, status.Error(codes.Unauthenticated, "missing auth header")
			}

			parts := strings.Split(authHeader[0], " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return nil, status.Error(codes.Unauthenticated, "invalid auth header")
			}
			token := parts[1]

			userId, err := a.jwtManager.ParseJWT(token)
			if err != nil {
				return nil, status.Error(codes.Unauthenticated, "invalid token")
			}

			ctx = context.WithValue(ctx, "userId", userId)
		}
		return handler(ctx, req)
	}
}
