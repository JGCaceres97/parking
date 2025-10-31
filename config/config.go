package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const (
	// HandlerTimeout es el tiempo máximo de espera para todas las operaciones del handler.
	HandlerTimeout = 60 * time.Second
	// DBTimeout es el tiempo máximo de espera para cualquier operación de DB.
	DBTimeout = 10 * time.Second
)

type Config struct {
	DBConnString  string
	JWTSecretKey  string
	ServerPort    string
	TokenDuration time.Duration
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Advertencia: No se pudo cargar el archivo .env. Usando variables de entorno o defaults.")
	}

	cfg := &Config{}

	dbHost := GetEnv("DB_HOST", "localhost")
	dbPort := GetEnv("DB_PORT", "3306")
	dbUser := GetEnv("MYSQL_USER", "root")
	dbPass := GetEnv("MYSQL_PASSWORD", "password")
	dbName := GetEnv("MYSQL_DATABASE", "parkingDb")

	cfg.DBConnString = fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=UTC",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbName,
	)

	cfg.JWTSecretKey = GetEnv("JWT_SECRET", "super_secreto_cambiar_1234")

	tokenDurationHours := GetEnv("TOKEN_DURATION_HOURS", "10")
	duration, err := time.ParseDuration(tokenDurationHours + "h")
	if err != nil {
		log.Printf("Advertencia: No se pudo parsear TOKEN_DURATION_HOURS. Usando 10h.")
		duration = 10 * time.Hour
	}
	cfg.TokenDuration = duration

	cfg.ServerPort = GetEnv("SERVER_PORT", "3000")

	return cfg
}

func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}
