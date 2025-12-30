package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/JGCaceres97/parking/internal/application/parking"
	"github.com/JGCaceres97/parking/internal/domain"
	"github.com/JGCaceres97/parking/internal/infrastructure/config"
)

type parkingRepository struct {
	DB *sql.DB
}

func NewParkingRepository(db *sql.DB) parking.Repository {
	return &parkingRepository{DB: db}
}

func (r *parkingRepository) CreateEntry(ctx context.Context, record *domain.ParkingRecord) error {
	ctx, cancel := context.WithTimeout(ctx, config.DBTimeout)
	defer cancel()

	query := `
		INSERT INTO PARKING_RECORDS
		(id, user_id, vehicle_type_id, license_plate, entry_time)
		VALUES (?, ?, ?, ?, ?);`

	_, err := r.DB.ExecContext(
		ctx,
		query,
		record.ID,
		record.UserID,
		record.VehicleTypeID,
		record.LicensePlate,
		record.EntryTime,
	)

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("timeout de DB excedido al buscar tipo de vehículo: %w", ctx.Err())
		}

		return fmt.Errorf("error al crear registro de entrada: %w", err)
	}

	return nil
}

func (r *parkingRepository) FindByID(ctx context.Context, id string) (*domain.ParkingRecord, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DBTimeout)
	defer cancel()

	query := `
		SELECT id, user_id, vehicle_type_id, license_plate, entry_time, exit_time, total_charge, calculated_hours
		FROM PARKING_RECORDS
		WHERE id = ?;`

	var record domain.ParkingRecord

	row := r.DB.QueryRowContext(ctx, query, id)

	var exitTime sql.NullTime
	var totalCharge sql.NullFloat64
	var calculatedHours sql.NullInt32

	err := row.Scan(
		&record.ID,
		&record.UserID,
		&record.VehicleTypeID,
		&record.LicensePlate,
		&record.EntryTime,
		&exitTime,
		&totalCharge,
		&calculatedHours,
	)

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("timeout de DB excedido al buscar registro de estacionamiento: %w", ctx.Err())
		}

		if err == sql.ErrNoRows {
			return nil, domain.ErrParkingRecordNotFound
		}

		return nil, fmt.Errorf("error al buscar registro de estacionamiento: %w", err)
	}

	if exitTime.Valid {
		record.ExitTime = &exitTime.Time
	}

	if totalCharge.Valid {
		record.TotalCharge = &totalCharge.Float64
	}

	if calculatedHours.Valid {
		h := int(calculatedHours.Int32)
		record.CalculatedHours = &h
	}

	return &record, nil
}

func (r *parkingRepository) FindOpenByLicensePlate(ctx context.Context, licensePlate string) (*domain.ParkingRecord, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DBTimeout)
	defer cancel()

	query := `
		SELECT id, user_id, vehicle_type_id, license_plate, entry_time
		FROM PARKING_RECORDS
		WHERE license_plate = ? AND exit_time IS NULL;`

	var record domain.ParkingRecord

	row := r.DB.QueryRowContext(ctx, query, licensePlate)

	err := row.Scan(
		&record.ID,
		&record.UserID,
		&record.VehicleTypeID,
		&record.LicensePlate,
		&record.EntryTime,
	)

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("timeout de DB excedido al buscar registro de estacionamiento abierto: %w", ctx.Err())
		}

		if err == sql.ErrNoRows {
			return nil, domain.ErrParkingRecordNotFound
		}

		return nil, fmt.Errorf("error al buscar registro de estacionamiento abierto: %w", err)
	}

	return &record, nil
}

func (r *parkingRepository) UpdateExit(ctx context.Context, record *domain.ParkingRecord) error {
	ctx, cancel := context.WithTimeout(ctx, config.DBTimeout)
	defer cancel()

	query := `
		UPDATE PARKING_RECORDS
		SET exit_time = ?, total_charge = ?, calculated_hours = ?
		WHERE id = ?;`

	result, err := r.DB.ExecContext(
		ctx,
		query,
		record.ExitTime,
		*record.TotalCharge,
		*record.CalculatedHours,
		record.ID,
	)

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("timeout de DB excedido al actualizar registro de salida: %w", ctx.Err())
		}

		return fmt.Errorf("error al actualizar registro de salida: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrParkingRecordNotFound
	}

	return nil
}

func (r *parkingRepository) ListCurrent(ctx context.Context) ([]domain.ParkingRecord, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DBTimeout)
	defer cancel()

	query := `
		SELECT id, user_id, vehicle_type_id, license_plate, entry_time
		FROM PARKING_RECORDS
		WHERE exit_time IS NULL
		ORDER BY entry_time DESC;`

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("timeout de DB excedido al listar vehículos actuales: %w", ctx.Err())
		}

		return nil, fmt.Errorf("error al ejecutar la consulta de vehículos actuales: %w", err)
	}
	defer rows.Close()

	records := []domain.ParkingRecord{}

	for rows.Next() {
		var record domain.ParkingRecord

		err := rows.Scan(
			&record.ID,
			&record.UserID,
			&record.VehicleTypeID,
			&record.LicensePlate,
			&record.EntryTime,
		)

		if err != nil {
			return nil, fmt.Errorf("error al escanear fila de registro actual: %w", err)
		}

		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar sobre resultados de registros actuales: %w", err)
	}

	return records, nil
}

func (r *parkingRepository) ListHistory(ctx context.Context) ([]domain.ParkingRecord, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DBTimeout)
	defer cancel()

	query := `
		SELECT id, user_id, vehicle_type_id, license_plate, entry_time, exit_time, total_charge, calculated_hours
		FROM PARKING_RECORDS
		WHERE exit_time IS NOT NULL
		ORDER BY exit_time DESC;`

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("timeout de DB excedido al listar historial: %w", ctx.Err())
		}

		return nil, fmt.Errorf("error al ejecutar la consulta de historial: %w", err)
	}
	defer rows.Close()

	var records []domain.ParkingRecord

	for rows.Next() {
		var record domain.ParkingRecord

		var exitTime sql.NullTime
		var totalCharge sql.NullFloat64
		var calculatedHours sql.NullInt32

		err := rows.Scan(
			&record.ID,
			&record.UserID,
			&record.VehicleTypeID,
			&record.LicensePlate,
			&record.EntryTime,
			&exitTime,
			&totalCharge,
			&calculatedHours,
		)

		if err != nil {
			return nil, fmt.Errorf("error al escanear fila de historial: %w", err)
		}

		if exitTime.Valid {
			record.ExitTime = &exitTime.Time
		}

		if totalCharge.Valid {
			record.TotalCharge = &totalCharge.Float64
		}

		if calculatedHours.Valid {
			h := int(calculatedHours.Int32)
			record.CalculatedHours = &h
		}

		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar sobre resultados de historial: %w", err)
	}

	return records, nil
}
