package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/melkomukovki/LockBox/api/pb"
	"github.com/melkomukovki/LockBox/internal/client/grpcclient"
	"github.com/melkomukovki/LockBox/internal/models"
)

// SecretService структура Secret сервиса
type SecretService struct {
	conn grpcclient.IGRPCClient
}

// List функция для получения списка секретов пользователя
func (s SecretService) List(ctx context.Context, token string) {
	fmt.Println("\tList of secrets")
	resp, err := s.conn.SecretsList(ctx, token)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	if len(resp.Secrets) == 0 {
		fmt.Println("No secrets found")
	}

	for _, secret := range resp.Secrets {
		fmt.Printf("ID: %d | Name: %s | Type: %s\n", secret.Id, secret.Name, secret.Type)
	}
}

// Get функция для получения конкретного секрета
func (s SecretService) Get(ctx context.Context, token string) {
	fmt.Println("\tGetting secret")
	id := inputInt("Secret ID > ")
	resp, err := s.conn.SecretsGet(ctx, token, id)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	secret := resp.Secret
	data := readSecretData(secret.Type, secret.Data)
	fmt.Printf("ID: %d\nName: %s\nType: %s\nDescription: %s\nData: \n%s\n",
		secret.Id, secret.Name, secret.Type, secret.Description, data)
}

// Delete функция для удаления конкретного секрета
func (s SecretService) Delete(ctx context.Context, token string) {
	fmt.Println("\tDeleting secret")
	id := inputInt("Secret ID > ")

	resp, err := s.conn.SecretsDelete(ctx, token, id)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println("Secret deleted", resp.Message)
}

// Add функция для добавления секрета
func (s SecretService) Add(ctx context.Context, token string) {
	fmt.Println("\tAdding secret")
	secret := inputSecret()
	resp, err := s.conn.SecretsAdd(ctx, token, secret)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Printf("Secret successfully saved! ID: %d\n", resp.Id)
}

// NewSecretService - функция конструктор для получения экземпляра Secret сервиса
func NewSecretService(conn grpcclient.IGRPCClient) ISecretService {
	if conn == nil {
		log.Fatal("GRPC client must not be nil")
	}
	return &SecretService{conn: conn}
}

func inputSecret() *pb.Secret {
	name := inputString("Secret Name > ")
	description := inputString("Secret Description > ")
	secretType := inputString("Secret Type > ")
	secretData := getSecret(secretType)

	return &pb.Secret{
		Name:        name,
		Type:        secretType,
		Description: description,
		Data:        secretData,
	}
}

func getSecret(t string) []byte {
	switch models.SecretType(t) {
	case models.Credentials:
		return getCredentialsSecret()
	case models.Text:
		return getTextSecret()
	case models.PaymentCard:
		return getPaymentCardSecret()
	case models.Binary:
		return getBinarySecret()
	default:
		log.Fatalf("Unknown secret type: %s", t)
		return nil
	}
}

func getCredentialsSecret() []byte {
	login := inputString("Secret Login > ")
	password := inputString("Secret Password > ")
	s := models.CredentialSecret{Login: login, Password: password}
	b, err := json.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func getTextSecret() []byte {
	input := inputString("Text to store > ")
	s := models.TextSecret{Text: input}
	b, err := json.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func getBinarySecret() []byte {
	input := inputString("Binary data to store > ")
	s := models.BinarySecret{Binary: input}
	b, err := json.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func getPaymentCardSecret() []byte {
	cardNumber := inputString("Payment Card Number > ")
	cardExpiryMonth := inputInt("Payment Card Expiry Month > ")
	cardExpiryYear := inputInt("Payment Card Expiry Year > ")
	cardCVV := inputString("Payment Card CVV > ")
	s := models.PaymentCardSecret{
		Number:      cardNumber,
		ExpiryMonth: int(cardExpiryMonth),
		ExpiryYear:  int(cardExpiryYear),
		CVV:         cardCVV,
	}
	b, err := json.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func readSecretData(t string, data []byte) string {
	switch models.SecretType(t) {
	case models.Text:
		var s models.TextSecret
		err := json.Unmarshal(data, &s)
		if err != nil {
			log.Fatal(err)
		}
		return s.Text
	case models.PaymentCard:
		var s models.PaymentCardSecret
		err := json.Unmarshal(data, &s)
		if err != nil {
			log.Fatal(err)
		}
		return fmt.Sprintf("Card Number: %s Until: %d/%d CVV: %s\n", s.Number, s.ExpiryMonth, s.ExpiryYear, s.CVV)
	case models.Binary:
		var s models.BinarySecret
		err := json.Unmarshal(data, &s)
		if err != nil {
			log.Fatal(err)
		}
		return fmt.Sprintf("Binary Data: %s", s.Binary)
	case models.Credentials:
		var s models.CredentialSecret
		err := json.Unmarshal(data, &s)
		if err != nil {
			log.Fatal(err)
		}
		return fmt.Sprintf("Login: %s\nPassword: %s\n", s.Login, s.Password)
	default:
		log.Fatalf("Unknown secret type: %s", t)
		return ""
	}

}
