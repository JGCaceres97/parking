package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/JGCaceres97/parking/config"
	"github.com/JGCaceres97/parking/internal/core/domain"
	"github.com/JGCaceres97/parking/internal/ports"
)

type InitRepository struct {
	DB *sql.DB
}

func NewInitRepository(db *sql.DB) ports.InitRepository {
	return &InitRepository{DB: db}
}

func (r *InitRepository) CreateAdmin(ctx context.Context, admin *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, config.DBTimeout)
	defer cancel()

	query := `
		INSERT INTO USERS (id, username, password_hash, role, is_active, created_at)
		VALUES (?, ?, ?, ?, ?, ?);`

	_, err := r.DB.ExecContext(
		ctx,
		query,
		admin.ID,
		admin.Username,
		admin.Password,
		admin.Role,
		admin.IsActive,
		admin.CreatedAt,
	)

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("timeout de DB excedido al crear usuario administrador: %w", ctx.Err())
		}

		return fmt.Errorf("error al crear usuario administrador: %w", err)
	}

	return nil
}

func (r *InitRepository) ExistsUsername(ctx context.Context, username string) bool {
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
