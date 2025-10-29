package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/JGCaceres97/parking/config"
	"github.com/JGCaceres97/parking/internal/core/domain"
	"github.com/JGCaceres97/parking/internal/ports"
)

type VehicleTypeRepository struct {
	DB *sql.DB
}

func NewVehicleTypeRepository(db *sql.DB) ports.VehicleTypeRepository {
	return &VehicleTypeRepository{DB: db}
}

func (r *VehicleTypeRepository) FindByID(ctx context.Context, id string) (*domain.VehicleType, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DBTimeout)
	defer cancel()

	query := `
		SELECT id, name, hourly_rate, description
		FROM VEHICLE_TYPES
		WHERE id = ?;`

	record := domain.VehicleType{}

	row := r.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&record.ID,
		&record.Name,
		&record.HourlyRate,
		&record.Description,
	)

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("timeout de DB excedido al buscar tipo de vehículo: %w", ctx.Err())
		}

		if err == sql.ErrNoRows {
			return nil, ports.ErrVehicleTypeNotFound
		}

		return nil, fmt.Errorf("error al buscar tipo de vehículo: %w", err)
	}

	return &record, nil
}

func (r *VehicleTypeRepository) ListAll(ctx context.Context) ([]domain.VehicleType, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DBTimeout)
	defer cancel()

	query := `
		SELECT id, name, hourly_rate, description
		FROM VEHICLE_TYPES
		ORDER BY name;`

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("timeout de DB excedido al listar tipos de vehículo: %w", ctx.Err())
		}

		return nil, fmt.Errorf("error al listar tipos de vehículo: %w", err)
	}
	defer rows.Close()

	vehicleTypes := []domain.VehicleType{}

	for rows.Next() {
		vt := domain.VehicleType{}

		err := rows.Scan(
			&vt.ID,
			&vt.Name,
			&vt.HourlyRate,
			&vt.Description,
		)

		if err != nil {
			return nil, fmt.Errorf("error al escanear fila de tipo de vehículo: %w", err)
		}

		vehicleTypes = append(vehicleTypes, vt)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar sobre resultados de tipos de vehículo: %w", err)
	}

	return vehicleTypes, nil
}
