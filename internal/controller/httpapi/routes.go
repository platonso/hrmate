package httpapi

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/platonso/hrmate/internal/domain"
	"github.com/platonso/hrmate/internal/service"
	"net/http"
)

type Router struct {
	handlerAuth *AuthHandler
	handlerForm *FormHandler
	handlerUser *UserHandler
	middleware  *AuthMiddleware
}

func NewRouter(
	authService *service.AuthService,
	userService *service.UserService,
	formService *service.FormService,
) *Router {
	v := validator.New()

	authMiddleware := &AuthMiddleware{
		authService: *authService,
		userService: *userService,
	}

	return &Router{
		handlerAuth: NewAuthHandler(authService, v),
		handlerUser: NewUsersHandler(userService, v),
		handlerForm: NewFormHandler(formService, v),
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
		).Group(func(r chi.Router) {
			r.Get("/users", rt.handlerUser.HandleGetUsers)
			r.Patch("/users/{id}/status", rt.handlerUser.HandleUpdateUserStatus)
		})
	})

	return r
}
