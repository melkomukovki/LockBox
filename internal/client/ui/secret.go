package ui

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"

	"github.com/melkomukovki/LockBox/api/pb"
	"github.com/melkomukovki/LockBox/internal/client/service"
	"github.com/melkomukovki/LockBox/internal/models"
)

type SecretHandler struct {
	srv *service.SecretService
}

func NewSecretHandler(srv *service.SecretService) *SecretHandler {
	return &SecretHandler{
		srv: srv,
	}
}

// List функция обработчик команд для получения списка секретов пользователя
func (sh *SecretHandler) List(ctx context.Context, token string) {
	displayHeader("List of secrets!")

	secrets, err := sh.srv.List(ctx, token)
	if err != nil {
		displayError(fmt.Errorf("Failed to retrieve secrets.\nCause: %v", err))
		return
	}

	rows := make([][]string, 0)
	for _, secret := range secrets {
		secretId := fmt.Sprintf("%d", secret.Id)
		rows = append(rows, []string{secretId, secret.Name, secret.Type, secret.Description})
	}

	t := table.New().
		Border(lipgloss.ThickBorder()).
		BorderStyle(tBorderStyle).
		StyleFunc(func(row, col int) lipgloss.Style {
			var style lipgloss.Style

			switch {
			case row == table.HeaderRow:
				return tHeaderStyle
			case row%2 == 0:
				style = tEvenRowStyle
			default:
				style = tOddRowStyle
			}

			if col == 0 {
				style = style.Width(5)
			}
			if col == 1 {
				style = style.Width(22)
			}
			if col == 3 {
				style = style.Width(33)
			}

			return style
		}).
		Headers("ID", "NAME", "TYPE", "DESCRIPTION").
		Rows(rows...)

	fmt.Println(t)
}

// Get функция обработчик команд для получения секрета по его id
func (sh *SecretHandler) Get(ctx context.Context, token string) {
	displayHeader("Getting secret!")

	id, err := inputInt("Secret ID > ")
	if err != nil {
		displayError(err)
		return
	}

	secret, err := sh.srv.Get(ctx, token, id)
	if err != nil {
		displayError(fmt.Errorf("Failed to retrieve secret.\nCause: %v", err))
		return
	}

	data, err := readSecretData(secret.Type, secret.Data)
	if err != nil {
		displayError(err)
		return
	}

	msg := fmt.Sprintf("\nID: %d\nName: %s\nType: %s\nDescription: %s\nData: \n%s",
		secret.Id, secret.Name, secret.Type, secret.Description, data)
	fmt.Println(outputStyle.Render(msg))
}

// Delete функция обработчик команд для удаления секрета
func (sh *SecretHandler) Delete(ctx context.Context, token string) {
	displayHeader("Deleting secret")

	id, err := inputInt("Secret ID > ")
	if err != nil {
		displayError(err)
		return
	}

	if err := sh.srv.Delete(ctx, token, id); err != nil {
		displayError(fmt.Errorf("Failed to delete secret.\nCause: %v", err))
		return
	}

	fmt.Println(successStyle.Render("Deleted!"))
}

// Create функция обработчик команд для создания секрета
func (sh *SecretHandler) Create(ctx context.Context, token string) {
	displayHeader("Creating secret")
	secret, err := inputSecret()
	if err != nil {
		displayError(err)
		return
	}

	secretId, err := sh.srv.Add(ctx, token, secret)
	if err != nil {
		displayError(fmt.Errorf("Failed to create secret.\nCause: %v\n", err))
		return
	}

	msg := fmt.Sprintf("Secret #%d created!", secretId)
	fmt.Println(successStyle.Render(msg))
}

func inputSecret() (*pb.Secret, error) {
	name, err := inputString("Secret Name > ")
	if err != nil {
		return nil, err
	}

	description, err := inputString("Secret Description > ")
	if err != nil {
		return nil, err
	}

	secretType, err := inputString("Secret Type > ")
	if err != nil {
		return nil, err
	}

	secretData, err := getSecret(secretType)
	if err != nil {
		return nil, err
	}

	return &pb.Secret{
		Name:        name,
		Type:        secretType,
		Description: description,
		Data:        secretData,
	}, nil
}

func getSecret(t string) ([]byte, error) {
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
		return nil, fmt.Errorf("Unknown secret type: %s", t)
	}
}

func getCredentialsSecret() ([]byte, error) {
	login, err := inputString("Secret Login > ")
	if err != nil {
		return nil, err
	}

	password, err := inputString("Secret Password > ")
	if err != nil {
		return nil, err
	}

	s := models.CredentialSecret{Login: login, Password: password}
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func getTextSecret() ([]byte, error) {
	input, err := inputString("Text to store > ")
	if err != nil {
		return nil, err
	}

	s := models.TextSecret{Text: input}
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func getBinarySecret() ([]byte, error) {
	input, err := inputString("Binary data to store > ")
	if err != nil {
		return nil, err
	}

	s := models.BinarySecret{Binary: input}
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func getPaymentCardSecret() ([]byte, error) {
	cardNumber, err := inputString("Payment Card Number > ")
	if err != nil {
		return nil, err
	}

	cardExpiryMonth, err := inputInt("Payment Card Expiry Month > ")
	if err != nil {
		return nil, err
	}

	cardExpiryYear, err := inputInt("Payment Card Expiry Year > ")
	if err != nil {
		return nil, err
	}

	cardCVV, err := inputString("Payment Card CVV > ")
	if err != nil {
		return nil, err
	}

	s := models.PaymentCardSecret{
		Number:      cardNumber,
		ExpiryMonth: int(cardExpiryMonth),
		ExpiryYear:  int(cardExpiryYear),
		CVV:         cardCVV,
	}
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func readSecretData(t string, data []byte) (string, error) {
	switch models.SecretType(t) {
	case models.Text:
		var s models.TextSecret
		err := json.Unmarshal(data, &s)
		if err != nil {
			return "", err
		}
		return s.Text, nil
	case models.PaymentCard:
		var s models.PaymentCardSecret
		err := json.Unmarshal(data, &s)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Card Number: %s Until: %d/%d CVV: %s\n", s.Number, s.ExpiryMonth, s.ExpiryYear, s.CVV), nil
	case models.Binary:
		var s models.BinarySecret
		err := json.Unmarshal(data, &s)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Binary Data: %s", s.Binary), nil
	case models.Credentials:
		var s models.CredentialSecret
		err := json.Unmarshal(data, &s)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Login: %s\nPassword: %s\n", s.Login, s.Password), nil
	default:
		return "", fmt.Errorf("unknown secret type: %s", t)
	}

}
