package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/JGCaceres97/parking/internal/adapters/api/handlers"
	"github.com/JGCaceres97/parking/internal/adapters/api/middlewares"
	"github.com/JGCaceres97/parking/internal/application/auth"
	"github.com/JGCaceres97/parking/internal/application/parking"
	"github.com/JGCaceres97/parking/internal/application/user"
	"github.com/JGCaceres97/parking/internal/application/vehicle_type"
	"github.com/JGCaceres97/parking/internal/domain"
	"github.com/JGCaceres97/parking/internal/infrastructure/config"
	"github.com/JGCaceres97/parking/web"
)

type routerConfig struct {
	auth        auth.Service
	parking     parking.Service
	user        user.Service
	vehicleType vehicle_type.Service
}

func New(auth auth.Service, parking parking.Service, user user.Service, vehicleType vehicle_type.Service) *routerConfig {
	return &routerConfig{
		auth,
		parking,
		user,
		vehicleType,
	}
}

func (rc *routerConfig) SetHandler() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(config.HandlerTimeout))

	authHandler := handlers.NewAuthHandler(rc.auth)
	parkingHandler := handlers.NewParkingHandler(rc.parking)
	userHandler := handlers.NewUserHandler(rc.user)
	vehicleTypeHandler := handlers.NewVehicleTypeHandler(rc.vehicleType)

	r.Handle("/*", web.Handler)

	r.Route("/api/v1", func(r chi.Router) {
		// Rutas públicas
		r.Post("/login", authHandler.Login)

		// Rutas protegidas
		r.Group(func(r chi.Router) {
			r.Use(middlewares.AuthMiddleware(rc.auth))

			// Parking
			r.Post("/parking/entry", parkingHandler.RecordEntry)
			r.Post("/parking/exit", parkingHandler.RecordExit)
			r.Get("/parking/{id}", parkingHandler.GetRecordByID)
			r.Get("/parking/current", parkingHandler.GetCurrentlyParked)
			r.Get("/parking/history", parkingHandler.GetHistory)

			// Users
			r.Put("/users/me", userHandler.UpdateUsername)

			// Vehicle types
			r.Get("/vehicle-types", vehicleTypeHandler.ListAll)

			// Admin
			r.Route("/admin", func(r chi.Router) {
				r.Use(middlewares.RoleMiddleware(domain.RoleAdmin))

				r.Get("/users", userHandler.ListUsers)
				r.Post("/users", userHandler.CreateUser)
				r.Put("/users/{userID}", userHandler.UpdateUser)
				r.Patch("/users/{userID}/active", userHandler.ToggleActiveStatus)
				r.Delete("/users/{userID}", userHandler.DeleteUser)
			})
		})
	})

	return r
}
