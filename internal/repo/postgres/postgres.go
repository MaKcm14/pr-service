package postgres

import (
	"log/slog"
)

// PostgreSQLRepo defines the logic of interaction with the PostgreSQL.
type PostgreSQLRepo struct {
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
		teamsRepo: teamsRepo{
			conf: conf,
		},
		usersRepo: usersRepo{
			conf: conf,
		},
		pullRequestRepo: pullRequestRepo{
			conf:      conf,
			usersRepo: usersRepo{conf},
		},
	}, nil
}
