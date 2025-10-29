package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/JGCaceres97/parking/config"
	"github.com/JGCaceres97/parking/internal/api"
	"github.com/JGCaceres97/parking/internal/api/handlers"
	"github.com/JGCaceres97/parking/internal/core/services"
	"github.com/JGCaceres97/parking/internal/infra/repository/mysql"
)

func main() {
	// Cargar configuración
	cfg := config.Load()

	// Inicializar conexión con DB
	db, err := initDB(cfg.DBConnString)
	if err != nil {
		log.Fatalf("error al inicializar la base de datos: %v", err)
	}
	defer func() {
		log.Println("Cerrando conexión a MySQL...")
		if err := db.Close(); err != nil {
			log.Printf("Advertencia: Error al cerrar la DB: %v", err)
		}

		log.Println("Conexión a MySQL cerrada.")
	}()
	log.Println("Conexión a MySQL establecida con éxito.")

	// Inyección de dependencias

	// -- A. Repositorios
	parkingRepo := mysql.NewParkingRepository(db)
	userRepo := mysql.NewUserRepository(db)
	vehicleTypeRepo := mysql.NewVehicleTypeRepository(db)

	// -- B. Servicios
	authService := services.NewAuthService(userRepo, cfg.JWTSecretKey, cfg.TokenDuration)
	parkingService := services.NewParkingService(parkingRepo, vehicleTypeRepo)
	userService := services.NewUserService(userRepo)
	vehicleTypeService := services.NewVehicleTypeService(vehicleTypeRepo)

	// -- C. Handlers
	authHandler := handlers.NewAuthHandler(authService)
	parkingHandler := handlers.NewParkingHandler(parkingService)
	userHandler := handlers.NewUserHandler(userService)
	vehicleTypeHandler := handlers.NewVehicleTypeHandler(vehicleTypeService)

	// Configuración del router
	routerCfg := api.RouterConfig{
		AuthService: authService,

		AuthHandler:        authHandler,
		ParkingHandler:     parkingHandler,
		UserHandler:        userHandler,
		VehicleTypeHandler: vehicleTypeHandler,
	}

	router := api.NewRouter(routerCfg)

	// Servidor
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	srv := &http.Server{Addr: addr, Handler: router}

	errors := make(chan error, 1)

	// Arrancar el servidor en una Go-routine.
	go func() {
		log.Printf("Servidor escuchando en %s", addr)

		errors <- srv.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Bloquear y esperar señal de apagado o error
	select {
	case err := <-errors:
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

func initDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, fmt.Errorf("error al abrir la conexión con la DB: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), config.DBTimeout)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error al hacer ping a la base de datos: %w", err)
	}

	return db, nil
}
