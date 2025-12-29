package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

func NewConnection(ctx context.Context, dsn string, timeout time.Duration) (*sql.DB, error) {
	conn, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("error al abrir la conexión con SQLite: %w", err)
	}

	conn.SetMaxOpenConns(1)
	conn.SetMaxIdleConns(1)

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err = conn.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error al hacer ping con SQLite: %w", err)
	}

	if _, err = conn.ExecContext(ctx, "PRAGMA optimize;"); err != nil {
		return nil, fmt.Errorf("error al ejecutar PRAGMA optimize: %w", err)
	}

	return conn, nil
}
