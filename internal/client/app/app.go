// Package app содержит код описывающий клиентское приложение и доступные команды для выполнения
package app

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/melkomukovki/LockBox/internal/client/ui"
)

// App структура приложения
type App struct {
	cmd           *cli.Command
	userService   *ui.UserHandler
	secretService *ui.SecretHandler
}

// NewApp конструктор для получения экземпляра приложения
func NewApp(userHandler *ui.UserHandler, secretHandler *ui.SecretHandler) *App {
	return &App{
		cmd: &cli.Command{
			Commands: []*cli.Command{
				{
					Name:    "user",
					Aliases: []string{"u"},
					Usage:   "option for user commands",
					Commands: []*cli.Command{
						{
							Name:  "register",
							Usage: "SignUp",
							Action: func(ctx context.Context, cmd *cli.Command) error {
								userHandler.SignUp(ctx)
								return nil
							},
						},
						{
							Name:  "login",
							Usage: "SignIn",
							Action: func(ctx context.Context, cmd *cli.Command) error {
								userHandler.SignIn(ctx)
								return nil
							},
						},
					},
				},
				{
					Name:    "secret",
					Aliases: []string{"s"},
					Usage:   "option for secret commands",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:    "token",
							Value:   "",
							Usage:   "user access token",
							Aliases: []string{"t"},
							Sources: cli.EnvVars("TOKEN"),
						},
					},
					Commands: []*cli.Command{
						{
							Name:  "list",
							Usage: "List of secrets",
							Action: func(ctx context.Context, cmd *cli.Command) error {
								secretHandler.List(ctx, cmd.String("token"))
								return nil
							},
						},
						{
							Name:  "get",
							Usage: "Get secret",
							Action: func(ctx context.Context, cmd *cli.Command) error {
								secretHandler.Get(ctx, cmd.String("token"))
								return nil
							},
						},
						{
							Name:  "delete",
							Usage: "Delete secret",
							Action: func(ctx context.Context, cmd *cli.Command) error {
								secretHandler.Delete(ctx, cmd.String("token"))
								return nil
							},
						},
						{
							Name:  "add",
							Usage: "Add secret",
							Action: func(ctx context.Context, cmd *cli.Command) error {
								secretHandler.Create(ctx, cmd.String("token"))
								return nil
							},
						},
					},
				},
			},
		},
	}
}

// Run - функция запуска приложения
func (a *App) Run(ctx context.Context, args []string) error {
	return a.cmd.Run(ctx, args)
}
