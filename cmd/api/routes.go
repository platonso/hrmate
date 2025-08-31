package main

import (
	"github.com/platonso/hrmate/internal/domain"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", app.HandleRegister(domain.RoleEmployee))
		r.Post("/login", app.HandleLogin)
	})

	r.Route("/hr/auth", func(r chi.Router) {
		r.Post("/register", app.HandleRegister(domain.RoleHR))
	})

	r.Route("/forms", func(r chi.Router) {
		r.With(
			app.AuthMiddleware,
			app.RequireRoles(domain.RoleEmployee),
			app.RequireActiveStatus,
		).Group(func(r chi.Router) {
			r.Post("/", app.HandleCreateForm)
			r.Get("/", app.HandleGetForms)
			r.Get("/{id}", app.HandleGetForm)
		})
	})

	r.Route("/hr", func(r chi.Router) {
		r.With(
			app.AuthMiddleware,
			app.RequireRoles(domain.RoleHR),
			app.RequireActiveStatus,
		).Group(func(r chi.Router) {
			r.Get("/users", app.HandleGetUsers)
			r.Get("/forms", app.HandleGetForms)
			r.Get("/forms/{id}", app.HandleGetForm)
			r.Patch("/forms/{id}/status", app.HandleUpdateFormStatus)
		})
	})

	r.Route("/admin", func(r chi.Router) {
		r.With(
			app.AuthMiddleware,
			app.RequireRoles(domain.RoleAdmin),
		).Group(func(r chi.Router) {
			r.Get("/users", app.HandleGetUsers)
			r.Patch("/users/{id}/status", app.HandleUpdateUserStatus)
		})
	})

	return r
}
