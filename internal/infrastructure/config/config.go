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
	AdminPassword string
	DBDriver      string
	DBConnString  string
	JWTSecretKey  string
	ServerPort    string
	TokenDuration time.Duration
}

func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Advertencia: No se pudo cargar el archivo .env. Usando variables de entorno o defaults.")
	}

	duration, err := time.ParseDuration(GetEnv("TOKEN_DURATION_HOURS", "10") + "h")
	if err != nil {
		log.Printf("Advertencia: No se pudo parsear TOKEN_DURATION_HOURS. Usando 10h.")
		duration = 10 * time.Hour
	}

	var dsn string
	driver := GetEnv("DB_DRIVER", "sqlite")

	switch driver {
	case "sqlite":
		dsn = GetEnv("SQLITE_DSN", "file:parking.db?_time_format=sqlite&_pragma=journal_mode(WAL)")

	case "mysql":
		dsn = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=UTC",
			GetEnv("MYSQL_USER", "root"),
			GetEnv("MYSQL_PASSWORD", "password"),
			GetEnv("DB_HOST", "localhost"),
			GetEnv("DB_PORT", "3306"),
			GetEnv("MYSQL_DATABASE", "parkingDb"),
		)

	default:
		log.Fatalf("unsupported DB_DRIVER: %s", driver)
	}

	return &Config{
		AdminPassword: GetEnv("ADMIN_PASSWORD", "admin"),
		DBDriver:      driver,
		DBConnString:  dsn,
		JWTSecretKey:  GetEnv("JWT_SECRET", "secret-key-to-sign-jwt"),
		ServerPort:    GetEnv("SERVER_PORT", "3000"),
		TokenDuration: duration,
	}
}
