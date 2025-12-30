package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JGCaceres97/parking/internal/adapters/api"
	"github.com/JGCaceres97/parking/internal/application/auth"
	"github.com/JGCaceres97/parking/internal/application/parking"
	"github.com/JGCaceres97/parking/internal/application/user"
	"github.com/JGCaceres97/parking/internal/application/vehicle_type"
	"github.com/JGCaceres97/parking/internal/domain"
	"github.com/JGCaceres97/parking/internal/infrastructure/config"
	"github.com/JGCaceres97/parking/internal/infrastructure/persistence"
	"github.com/JGCaceres97/parking/pkg/ulid"
)

func main() {
	cfg := config.Load()

	db, err := persistence.NewConnection(
		context.Background(),
		cfg.DBDriver,
		cfg.DBConnString,
		config.DBTimeout)

	if err != nil {
		log.Fatalf("error al inicializar la conexión con base de datos: %v", err)
	}

	defer func() {
		log.Println("Cerrando conexión a DB...")
		if err := db.Close(); err != nil {
			log.Printf("Advertencia: Error al cerrar la DB: %v", err)
		}

		log.Println("Conexión a DB cerrada.")
	}()

	// Inyección de dependencias
	// -- A. Repositorios
	repos := persistence.NewRepositories(db, cfg.DBDriver)

	// -- B. Servicios
	authService := auth.NewService(repos.User, cfg.JWTSecretKey, cfg.TokenDuration)
	parkingService := parking.NewService(repos.Parking, repos.VehicleType)
	userService := user.NewService(repos.User)
	vehicleTypeService := vehicle_type.NewService(repos.VehicleType)

	// Admin User
	if err := ensureAdminUser(context.Background(), userService, cfg.AdminPassword); err != nil {
		log.Fatalf("error asegurando usuario administrador: %v", err)
	}

	// Configuración del router
	handler := api.New(authService, parkingService, userService, vehicleTypeService).SetHandler()

	// Servidor
	srv := &http.Server{Addr: ":" + cfg.ServerPort, Handler: handler}
	start(srv)
}

func ensureAdminUser(ctx context.Context, service user.Service, password string) error {
	admin := &domain.User{
		ID:        ulid.GenerateNewULID(),
		Username:  domain.AdminUsername,
		Password:  password,
		Role:      domain.RoleAdmin,
		IsActive:  true,
		CreatedAt: time.Now().UTC().Truncate(time.Second),
	}

	_, err := service.Create(ctx, admin)
	if err != nil && !errors.Is(err, domain.ErrUsernameAlreadyExists) {
		return err
	}

	return nil
}

func start(srv *http.Server) {
	errCh := make(chan error, 1)

	// Arrancar el servidor en una Go-routine.
	go func() {
		log.Printf("Servidor escuchando en %s", srv.Addr)

		errCh <- srv.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Bloquear y esperar señal de apagado o error
	select {
	case err := <-errCh:
		log.Fatalf("Error de servidor: %v", err)
	case sig := <-quit:
		log.Printf("Recibida señal '%v'. Iniciando proceso...", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("El servidor se apagó forzosamente: %v", err)
		}

		log.Println("Servidor detenido con éxito.")
	}
}
