// Package service - содержит сервисный слов нашего клиента
// Все сервисы работают через получения пользовательского ввода
package service

import (
	"context"
)

// IUserService - интерфейс описывающий методы, которые должен реализовать User сервис
type IUserService interface {
	Login(ctx context.Context)
	Register(ctx context.Context)
}

// ISecretService - интерфейс описывающий методы, которые должен реализовать Secret сервис
type ISecretService interface {
	List(ctx context.Context, token string)
	Get(ctx context.Context, token string)
	Delete(ctx context.Context, token string)
	Add(ctx context.Context, token string)
}
