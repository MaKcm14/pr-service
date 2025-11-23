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

func newPostgresConfig(log *slog.Logger, dsn string) (*postgresConfig, error) {
	const op = "postgres.new-config"

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	conn, err := pgxpool.New(ctx, dsn)
	if err != nil {
		retErr := fmt.Errorf("error of the %s: %w: %w", op, repo.ErrConnToRepository, err)
		log.Error(retErr.Error())
		return nil, retErr
	}

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

func getSqlViewBool(val bool) string {
	if val {
		return "TRUE"
	}
	return "FALSE"
}

func (p postgresConfig) close() {
	p.conn.Close()
}
