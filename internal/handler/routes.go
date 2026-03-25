package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
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

	// Authentication
	r.Post("/auth/register", rt.handlerAuth.HandleRegister)
	r.Post("/auth/login", rt.handlerAuth.HandleLogin)

	// Forms (Employee)
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

	// HR routes
	r.Route("/hr", func(r chi.Router) {
		r.With(
			rt.middleware.AuthMiddleware,
			rt.middleware.RequireRoles(domain.RoleHR),
			rt.middleware.RequireActiveStatus,
		).Group(func(r chi.Router) {
			r.Get("/users", rt.handlerUser.HandleGetUsers)
			r.Get("/users/{id}/forms", rt.handlerForm.HandleGetFormsWithUser)

			r.Get("/forms", rt.handlerForm.HandleGetFormsWithUsers)
			r.Get("/forms/{id}", rt.handlerForm.HandleGetForm)
			r.Patch("/forms/{id}/approve", rt.handlerForm.HandleApprove)
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
			r.Patch("/users/{id}/status", rt.handlerUser.HandleUpdateUserStatus)
		})
	})

	return r
}
