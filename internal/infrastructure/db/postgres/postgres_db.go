package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type postgres struct {
	DB *pgx.Conn
}

func NewWithConn(addr string) (*postgres, error) {
	ctx := context.Background()
	db, err := pgx.Connect(ctx, addr)
	if err != nil {
		return nil, fmt.Errorf("NewWithConn: failed connect to postgresql: %w", err)
	}

	if err := db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("NewWithConn: failed ping database: %w", err)
	}

	return &postgres{
		DB: db,
	}, nil
}

func (p *postgres) CloseConn(ctx context.Context) error {
	if err := p.DB.Close(ctx); err != nil {
		return fmt.Errorf("CloseConn: failed close conn postgres: %w", err)
	}

	return nil
}
