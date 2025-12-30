package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/JGCaceres97/parking/internal/application/parking"
	"github.com/JGCaceres97/parking/internal/application/user"
	"github.com/JGCaceres97/parking/internal/application/vehicle_type"
	"github.com/JGCaceres97/parking/internal/infrastructure/persistence/mysql"
	"github.com/JGCaceres97/parking/internal/infrastructure/persistence/sqlite"
)

type repositories struct {
	Parking     parking.Repository
	User        user.Repository
	VehicleType vehicle_type.Repository
}

func NewConnection(ctx context.Context, driver, dsn string, timeout time.Duration) (*sql.DB, error) {
	switch driver {
	case "sqlite":
		return sqlite.NewConnection(ctx, dsn, timeout)

	case "mysql":
		return mysql.NewConnection(ctx, dsn, timeout)

	default:
		return nil, fmt.Errorf("driver de DB no soportado: %s", driver)
	}
}

func NewRepositories(db *sql.DB, driver string) *repositories {
	switch driver {
	case "sqlite", "mysql":
		return &repositories{
			Parking:     mysql.NewParkingRepository(db),
			User:        mysql.NewUserRepository(db),
			VehicleType: mysql.NewVehicleTypeRepository(db),
		}

	default:
		panic("driver de DB no soportado: " + driver)
	}
}
