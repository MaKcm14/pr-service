package postgres

import (
	"log/slog"
)

// PostgreSQLRepo defines the logic of interaction with the PostgreSQL.
type PostgreSQLRepo struct {
	conf *postgresConfig

	teamsRepo
	usersRepo
	pullRequestRepo
}

func New(log *slog.Logger, socket string) (PostgreSQLRepo, error) {
	conf, err := newPostgresConfig(log, socket)
	if err != nil {
		return PostgreSQLRepo{}, err
	}
	return PostgreSQLRepo{
		conf: conf,
		teamsRepo: teamsRepo{
			conf: conf,
		},
		usersRepo: usersRepo{
			conf: conf,
		},
		pullRequestRepo: pullRequestRepo{
			conf: conf,
		},
	}, nil
}
