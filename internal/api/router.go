package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/JGCaceres97/parking/internal/api/handlers"
	"github.com/JGCaceres97/parking/internal/api/middlewares"
	"github.com/JGCaceres97/parking/internal/core/domain"
	"github.com/JGCaceres97/parking/internal/ports"
	"github.com/JGCaceres97/parking/web"
)

type RouterConfig struct {
	AuthService ports.AuthService

	AuthHandler        *handlers.AuthHandler
	ParkingHandler     *handlers.ParkingHandler
	UserHandler        *handlers.UserHandler
	VehicleTypeHandler *handlers.VehicleTypeHandler
}

func NewRouter(cfg RouterConfig) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Handle("/*", web.Handler)

	r.Route("/api/v1", func(r chi.Router) {
		// Rutas públicas
		r.Post("/login", cfg.AuthHandler.Login)

		// Rutas protegidas
		r.Group(func(r chi.Router) {
			r.Use(middlewares.AuthMiddleware(cfg.AuthService))

			// Parking
			r.Post("/parking/entry", cfg.ParkingHandler.RecordEntry)
			r.Post("/parking/exit", cfg.ParkingHandler.RecordExit)
			r.Get("/parking/{id}", cfg.ParkingHandler.GetRecordByID)
			r.Get("/parking/current", cfg.ParkingHandler.GetCurrentlyParked)
			r.Get("/parking/history", cfg.ParkingHandler.GetHistory)

			// Users
			r.Put("/users/me", cfg.UserHandler.UpdateMyProfile)

			// Vehicle types
			r.Get("/vehicle-types", cfg.VehicleTypeHandler.ListAll)

			// Admin
			r.Route("/admin", func(r chi.Router) {
				r.Use(middlewares.RoleMiddleware(domain.RoleAdmin))

				r.Get("/users", cfg.UserHandler.ListUsers)
				r.Post("/users", cfg.UserHandler.CreateUser)
				r.Put("/users/{userID}", cfg.UserHandler.UpdateUser)
				r.Patch("/users/{userID}/active", cfg.UserHandler.ToggleActiveStatus)
				r.Delete("/users/{userID}", cfg.UserHandler.DeleteUser)
			})
		})
	})

	return r
}
