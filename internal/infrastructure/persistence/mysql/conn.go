package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func NewConnection(ctx context.Context, dsn string, timeout time.Duration) (*sql.DB, error) {
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error al abrir la conexión con MySQL: %w", err)
	}

	conn.SetMaxOpenConns(20)
	conn.SetMaxIdleConns(20)
	conn.SetConnMaxLifetime(2 * time.Minute)

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err = conn.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error al hacer ping con MySQL: %w", err)
	}

	return conn, nil
}
