package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/platonso/hrmate/internal/config"
	"github.com/platonso/hrmate/internal/handler"
	"github.com/platonso/hrmate/internal/repository/postgres"
	"github.com/platonso/hrmate/internal/service/auth"
	"github.com/platonso/hrmate/internal/service/form"
	"github.com/platonso/hrmate/internal/service/user"
)

type Application struct {
	Config *config.Config
	Auth   *auth.Service
	Users  *user.Service
	Forms  *form.Service

	router    *handler.Router
	closeFunc func()
}

func New(ctx context.Context, cfg *config.Config) (*Application, error) {
	postgresRepo, err := postgres.NewRepository(ctx, cfg.GetConnStr())
	if err != nil {
		return nil, fmt.Errorf("failed to create repository: %w", err)
	}

	userService := user.NewService(&postgresRepo.Users)
	authService := auth.NewService(&postgresRepo.Users)
	formService := form.NewService(&postgresRepo.Forms, &postgresRepo.Users)

	if err := authService.ImplementAdmin(ctx); err != nil {
		return nil, fmt.Errorf("failed to implement admin: %w", err)
	}

	router := handler.NewRouter(authService, userService, formService)

	app := &Application{
		Config: cfg,
		Auth:   authService,
		Users:  userService,
		Forms:  formService,

		router:    router,
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
