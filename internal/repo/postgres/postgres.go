package postgres

import (
	"log/slog"
	"sync"
)

// PostgreSQLRepo defines the logic of interaction with the PostgreSQL.
type PostgreSQLRepo struct {
	conf *postgresConfig

	teamsRepo teamsRepo
	usersRepo usersRepo
	prRepo    pullRequestRepo

	closer sync.Once
}

func New(log *slog.Logger, socket string) (*PostgreSQLRepo, error) {
	conf, err := newPostgresConfig(log, socket)
	if err != nil {
		return &PostgreSQLRepo{}, err
	}
	return &PostgreSQLRepo{
		conf: conf,
		teamsRepo: teamsRepo{
			conf: conf,
		},
		usersRepo: usersRepo{
			conf: conf,
		},
		prRepo: pullRequestRepo{
			conf: conf,
		},
	}, nil
}

func (p *PostgreSQLRepo) Close() {
	p.closer.Do(func() {
		p.conf.close()
	})
}
