package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/JGCaceres97/parking/internal/application/user"
	"github.com/JGCaceres97/parking/internal/domain"
	"github.com/JGCaceres97/parking/internal/infrastructure/config"
)

type userRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) user.Repository {
	return &userRepository{DB: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, config.DBTimeout)
	defer cancel()

	query := `
		INSERT INTO USERS (id, username, password_hash, role, is_active, created_at)
		VALUES (?, ?, ?, ?, ?, ?);`

	_, err := r.DB.ExecContext(
		ctx,
		query,
		user.ID,
		strings.ToLower(user.Username),
		user.Password,
		user.Role,
		user.IsActive,
		user.CreatedAt,
	)

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("timeout de DB excedido al crear usuario: %w", ctx.Err())
		}

		return fmt.Errorf("error al crear usuario: %w", err)
	}

	return nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DBTimeout)
	defer cancel()

	query := `
		SELECT id, username, role, is_active, created_at
		FROM USERS
		WHERE id = ?;`

	var user domain.User

	row := r.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
	)

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("timeout de DB excedido al buscar usuario: %w", ctx.Err())
		}

		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}

		return nil, fmt.Errorf("error al buscar usuario: %w", err)
	}

	return &user, nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DBTimeout)
	defer cancel()

	query := `
		SELECT id, username, password_hash, role, is_active, created_at
		FROM USERS
		WHERE username = ?;`

	var user domain.User

	row := r.DB.QueryRowContext(ctx, query, username)

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
	)

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("timeout de DB excedido al buscar usuario: %w", ctx.Err())
		}

		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}

		return nil, fmt.Errorf("error al buscar usuario: %w", err)
	}

	return &user, nil
}

func (r *userRepository) ExistsUsername(ctx context.Context, username string) bool {
	ctx, cancel := context.WithTimeout(ctx, config.DBTimeout)
	defer cancel()

	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM USERS WHERE username = ?);"

	err := r.DB.QueryRowContext(ctx, query, username).Scan(&exists)
	if err != nil {
		return false
	}

	return exists
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, config.DBTimeout)
	defer cancel()

	var exists bool
	checkQuery := "SELECT EXISTS(SELECT 1 FROM USERS WHERE id = ?);"

	err := r.DB.QueryRowContext(ctx, checkQuery, user.ID).Scan(&exists)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("timeout de DB excedido al verificar existencia de usuario: %w", err)
		}

		return fmt.Errorf("error al verificar existencia de usuario: %w", err)
	}

	if !exists {
		return domain.ErrUserNotFound
	}

	updateQuery := `
		UPDATE USERS
		SET username = ?, role = ?, is_active = ?
		WHERE id = ?;`

	result, err := r.DB.ExecContext(
		ctx,
		updateQuery,
		user.Username,
		user.Role,
		user.IsActive,
		user.ID,
	)

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("timeout de DB excedido al actualizar usuario: %w", err)
		}

		return fmt.Errorf("error al actualizar usuario: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, config.DBTimeout)
	defer cancel()

	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM USERS WHERE id = ?);`

	err := r.DB.QueryRowContext(ctx, checkQuery, id).Scan(&exists)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("timeout de DB excedido al verificar existencia de usuario: %w", err)
		}

		return fmt.Errorf("error al verificar existencia de usuario: %w", err)
	}

	if !exists {
		return domain.ErrUserNotFound
	}

	deleteQuery := `
		DELETE FROM USERS
		WHERE id = ?;`

	result, err := r.DB.ExecContext(ctx, deleteQuery, id)

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("timeout de DB excedido al eliminar usuario: %w", err)
		}

		return fmt.Errorf("error al eliminar usuario: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil
	}

	return nil
}

func (r *userRepository) ListAll(ctx context.Context, id string) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DBTimeout)
	defer cancel()

	query := `
		SELECT id, username, role, is_active, created_at
		FROM USERS
		WHERE id != ?;`

	rows, err := r.DB.QueryContext(ctx, query, id)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("timeout de DB excedido al listar usuarios: %w", ctx.Err())
		}

		return nil, fmt.Errorf("error al ejecutar la consulta de listado de usuarios: %w", err)
	}
	defer rows.Close()

	users := []domain.User{}

	for rows.Next() {
		var user domain.User

		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Role,
			&user.IsActive,
			&user.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("error al escanear fila de usuario: %w", err)
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar sobre resultados de usuarios: %w", err)
	}

	return users, nil
}
