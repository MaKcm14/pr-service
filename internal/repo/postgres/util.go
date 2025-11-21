package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/MaKcm14/pr-service/internal/repo"
	"github.com/jackc/pgx/v5/pgxpool"
)

// postgresConfig defines the PostgreSQL repo configuration object.
type postgresConfig struct {
	log  *slog.Logger
	conn *pgxpool.Pool
}

func newPostgresConfig(log *slog.Logger, socket string) (*postgresConfig, error) {
	const op = "postgres.new-config"

	conn, err := pgxpool.New(context.Background(), socket)
	if err != nil {
		retErr := fmt.Errorf("error of the %s: %w: %w", op, repo.ErrConnToRepository, err)
		log.Error(retErr.Error())
		return nil, retErr
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	if err := conn.Ping(ctx); err != nil {
		retErr := fmt.Errorf("error of the %s: %w: %w", op, repo.ErrConnToRepository, err)
		log.Error(retErr.Error())
		return nil, retErr
	}

	return &postgresConfig{
		log:  log,
		conn: conn,
	}, nil
}
