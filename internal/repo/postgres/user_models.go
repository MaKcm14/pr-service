package postgres

import (
	"context"
	"fmt"

	"github.com/MaKcm14/pr-service/internal/entities"
	"github.com/MaKcm14/pr-service/internal/repo"
)

// usersRepo defines the logic of interaction with the users models
type usersRepo struct {
	conf *postgresConfig
}

const updateUser = `
	UPDATE users
	SET is_active=$1
	WHERE id=$2
`

func (u usersRepo) SetUserIsActive(ctx context.Context, dto entities.User) (entities.User, error) {
	const op = "postgres.set-user-is-active"

	tag, err := u.conf.conn.Exec(ctx, updateUser, true, dto.ID)
	if err != nil {
		retErr := fmt.Errorf("error of the %s: %w: %w", op, repo.ErrQueryExec, err)
		u.conf.log.Warn(retErr.Error())
		return entities.User{}, retErr
	}

	if tag.RowsAffected() == 0 {
		return entities.User{}, repo.ErrModelNotFound
	}

	user, err := u.getUser(ctx, dto.ID)
	if err != nil {
		return entities.User{}, err
	}

	return user, nil
}

const selectUserTeamName = `
	SELECT users.id, users.username, users.is_active, teams.team_name
	FROM users
	JOIN teams ON users.team_id=teams.id
	WHERE users.id=$1
`

func (u usersRepo) getUser(ctx context.Context, id entities.UserID) (entities.User, error) {
	const op = "postgres.get-user-team-name"

	rows, err := u.conf.conn.Query(ctx, selectUserTeamName, id)
	if err != nil {
		retErr := fmt.Errorf("error of the %s: %w: %s", op, repo.ErrQueryExec, err)
		u.conf.log.Warn(retErr.Error())
		return entities.User{}, retErr
	}
	defer rows.Close()

	user := entities.User{}
	if rows.Next() {
		rows.Scan(&user.ID, &user.Name, &user.IsActive, &user.TeamName)
		return user, nil
	} else if rows.Err() != nil {
		retErr := fmt.Errorf("error of the %s: %w: %s", op, repo.ErrResProcessing, err)
		u.conf.log.Warn(retErr.Error())
		return entities.User{}, retErr
	}

	return entities.User{}, repo.ErrModelNotFound
}
