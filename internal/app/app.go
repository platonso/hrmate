package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/platonso/hrmate/internal/config"
	"github.com/platonso/hrmate/internal/controller/httpapi"
	"github.com/platonso/hrmate/internal/repository/postgres"
	"github.com/platonso/hrmate/internal/service"
	"log"
	"net/http"
)

type Application struct {
	Config *config.Config
	Auth   *service.AuthService
	Users  *service.UserService
	Forms  *service.FormService

	router    *httpapi.Router
	closeFunc func()
}

func New(ctx context.Context, cfg *config.Config) (*Application, error) {
	postgresRepo, err := postgres.NewRepository(ctx, cfg.GetConnStr())
	if err != nil {
		return nil, fmt.Errorf("failed to create repository: %w", err)
	}

	authService := service.NewAuthService(postgresRepo.Users)
	userService := service.NewUserService(postgresRepo.Users)
	formService := service.NewFormService(postgresRepo.Forms, postgresRepo.Users)

	app := &Application{
		Config: cfg,
		Auth:   authService,
		Users:  userService,
		Forms:  formService,

		router:    httpapi.NewRouter(authService, userService, formService),
		closeFunc: postgresRepo.Close,
	}

	return app, nil
}

func (app *Application) Run() error {
	return app.StartServer()
}

func (app *Application) routes() http.Handler {
	return app.router.Routes()
}

func (app *Application) StartServer() error {
	server := &http.Server{
		Addr:    ":" + app.Config.HTTPPort,
		Handler: app.routes(),
	}

	log.Printf("Starting server on port %s", app.Config.HTTPPort)
	if err := server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
	return nil
}

func (app *Application) Close() {
	if app.closeFunc != nil {
		app.closeFunc()
	}
}
