package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	"github.com/platonso/hrmate/internal/domain"
	"github.com/platonso/hrmate/internal/handler/auth"
	"github.com/platonso/hrmate/internal/handler/form"
	"github.com/platonso/hrmate/internal/handler/middleware"
	"github.com/platonso/hrmate/internal/handler/user"
)

type AuthProvider interface {
	auth.Service
	middleware.AuthService
}

type UserProvider interface {
	user.Service
	middleware.UserService
}

type Router struct {
	handlerAuth *auth.Handler
	handlerUser *user.Handler
	handlerForm *form.Handler
	middleware  *middleware.Auth
}

func NewRouter(authSvc AuthProvider, userSvc UserProvider, formSvc form.Service,
) *Router {
	v := validator.New()

	authMiddleware := &middleware.Auth{
		AuthSvc: authSvc,
		UserSvc: userSvc,
	}

	return &Router{
		handlerAuth: auth.NewHandler(authSvc, v),
		handlerUser: user.NewHandler(userSvc, v),
		handlerForm: form.NewHandler(formSvc, v),
		middleware:  authMiddleware,
	}
}

func (rt *Router) Routes() http.Handler {
	r := chi.NewRouter()

	// CORS middleware ==================================================================
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins for testing
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	// ===================================================================================

	// Authentication
	r.Post("/register", rt.handlerAuth.HandleRegister)
	r.Post("/login", rt.handlerAuth.HandleLogin)

	// Employee
	r.Route("/forms", func(r chi.Router) {
		r.With(
			rt.middleware.AuthMiddleware,
			rt.middleware.RequireRoles(domain.RoleEmployee),
			rt.middleware.RequireActiveStatus,
		).Group(func(r chi.Router) {
			r.Post("/", rt.handlerForm.HandleCreateForm)
			r.Get("/", rt.handlerForm.HandleGetForms)
			r.Get("/{id}", rt.handlerForm.HandleGetForm)
		})
	})

	// HR
	r.Route("/hr", func(r chi.Router) {
		r.With(
			rt.middleware.AuthMiddleware,
			rt.middleware.RequireRoles(domain.RoleHR),
			rt.middleware.RequireActiveStatus,
		).Group(func(r chi.Router) {
			r.Get("/users", rt.handlerUser.HandleGetUsers)

			r.Get("/forms", rt.handlerForm.HandleGetFormsWithUsers)
			r.Get("/forms/{id}", rt.handlerForm.HandleGetForm)
			r.Patch("/forms/{id}/approve", rt.handlerForm.HandleApprove)
			r.Patch("/forms/{id}/reject", rt.handlerForm.HandleReject)
		})
	})

	// Administration
	r.Route("/admin", func(r chi.Router) {
		r.With(
			rt.middleware.AuthMiddleware,
			rt.middleware.RequireRoles(domain.RoleAdmin),
			rt.middleware.RequireActiveStatus,
		).Group(func(r chi.Router) {
			r.Get("/users", rt.handlerUser.HandleGetUsers)
			r.Patch("/users/{id}/activate", rt.handlerUser.HandleActivate)
			r.Patch("/users/{id}/deactivate", rt.handlerUser.HandleDeactivate)
		})
	})

	return r
}
