package ui

import (
	"context"
	"fmt"

	"github.com/melkomukovki/LockBox/internal/client/service"
)

type UserHandler struct {
	srv *service.UserService
}

func NewUserHandler(srv *service.UserService) *UserHandler {
	return &UserHandler{
		srv: srv,
	}
}

func (uh *UserHandler) SignUp(ctx context.Context) {
	displayHeader("Sign Up!")

	login, err := inputString("Login > ")
	if err != nil {
		displayError(err)
		return
	}

	password, err := inputString("Password > ")
	if err != nil {
		displayError(err)
		return
	}

	confirmPassword, err := inputString("ConfirmPassword > ")
	if err != nil {
		displayError(err)
		return
	}

	if password != confirmPassword {
		fmt.Println(errorStyle.Render("Passwords do not match!"))
		return
	}

	result, err := uh.srv.Register(ctx, login, password)
	if err != nil {
		displayError(err)
		return
	}

	msg := fmt.Sprintf("Status: %s", result)
	fmt.Println(successStyle.Render(msg))
}

func (uh *UserHandler) SignIn(ctx context.Context) {
	displayHeader("Sign In!")
	login, err := inputString("Login > ")
	if err != nil {
		displayError(err)
		return
	}

	password, err := inputString("Password > ")
	if err != nil {
		displayError(err)
		return
	}

	result, err := uh.srv.Login(ctx, login, password)
	if err != nil {
		displayError(fmt.Errorf("Failed to login.\nCause: %v", err))
		return
	}
	fmt.Println(successStyle.Render("Success!"))
	fmt.Printf(outputStyle.Render("Access Token: %s"), result)
}
