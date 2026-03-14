package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/platonso/hrmate/internal/domain"
	auth2 "github.com/platonso/hrmate/internal/handler/auth"
	form2 "github.com/platonso/hrmate/internal/handler/form"
	"github.com/platonso/hrmate/internal/handler/httpapi"
	user2 "github.com/platonso/hrmate/internal/handler/user"
	"github.com/platonso/hrmate/internal/service/auth"
	"github.com/platonso/hrmate/internal/service/form"
	"github.com/platonso/hrmate/internal/service/user"
)

type Router struct {
	handlerAuth *auth2.Handler
	handlerForm *form2.Handler
	handlerUser *user2.Handler
	middleware  *httpapi.AuthMiddleware
}

func NewRouter(
	authService *auth.Service,
	userService *user.Service,
	formService *form.Service,
) *Router {
	v := validator.New()

	authMiddleware := &httpapi.AuthMiddleware{
		AuthSvc: authService,
		UserSvc: userService,
	}

	return &Router{
		handlerAuth: auth2.NewHandler(authService, v),
		handlerUser: user2.NewHandler(userService, v),
		handlerForm: form2.NewHandler(formService, v),
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
			r.Get("/forms", rt.handlerForm.HandleGetForms)
			r.Get("/forms/{id}", rt.handlerForm.HandleGetForm)
			r.Patch("/forms/{id}/status", rt.handlerForm.HandleUpdateFormStatus)
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
