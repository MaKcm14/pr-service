package postgres

import (
	"bytes"
	"context"
	"fmt"

	"github.com/MaKcm14/pr-service/internal/entities"
	"github.com/MaKcm14/pr-service/internal/repo"
)

// teamsRepo defines the logic of interaction with the teams models.
type teamsRepo struct {
	conf *postgresConfig
}

// GetTeam defines the logic of getting the team object for the current name.
func (t teamsRepo) GetTeam(ctx context.Context, name string) (entities.Team, bool, error) {
	const op = "postgres.get-team"

	team, ok, err := t.isTeamExists(ctx, name)
	if err != nil {
		retErr := fmt.Errorf("error of the %s: %w", op, err)
		t.conf.log.Warn(retErr.Error())
		return entities.Team{}, false, retErr
	} else if !ok {
		return entities.Team{}, false, nil
	}

	members, err := t.getTeamMembers(ctx, name)
	if err != nil {
		retErr := fmt.Errorf("error of the %s: %w", op, err)
		t.conf.log.Warn(retErr.Error())
		return entities.Team{}, false, retErr
	}
	team.Members = members

	return team, true, nil
}

// CreateTeam defines the logic of creating a new team.
func (t teamsRepo) CreateTeam(ctx context.Context, team entities.Team) error {
	const op = "postgres.create-team"

	team, ok, err := t.isTeamExists(ctx, team.Name)
	if err != nil {
		retErr := fmt.Errorf("error of the %s: %w", op, err)
		t.conf.log.Warn(retErr.Error())
		return retErr
	} else if ok {
		return repo.ErrCreateMultipleUniqueModels
	}

	if err := t.createTeam(ctx, team); err != nil {
		return err
	}
	return nil
}

const insertTeam = `
	INSERT INTO teams (team_name)
	VALUES ($1)
`

// createTeam defines the logic of creating the 'Team' model.
func (t teamsRepo) createTeam(ctx context.Context, team entities.Team) error {
	const op = "postgres.create-team"

	_, err := t.conf.conn.Exec(ctx, insertTeam, team.Name)

	if err != nil {
		retErr := fmt.Errorf("error of the %s: %w: %w", op, repo.ErrQueryExec, err)
		t.conf.log.Warn(retErr.Error())
		return retErr
	}

	if err := t.addMembersList(ctx, team); err != nil {
		return err
	}
	return nil
}

const insertMembers = `
INSERT INTO users (id, username, is_active)
`

// addMembersList defines the logic of adding the members list for the current team.
func (t teamsRepo) addMembersList(ctx context.Context, team entities.Team) error {
	const op = "postgres.add-members-list"

	if len(team.Members) == 0 {
		return nil
	}

	query := bytes.Buffer{}
	query.WriteString(insertMembers)
	for _, user := range team.Members {
		query.WriteString(fmt.Sprintf("(%d, %s, %v)\n", user.ID, user.Name, user.IsActive))
	}

	_, err := t.conf.conn.Exec(ctx, query.String())
	if err != nil {
		retErr := fmt.Errorf("error of the %s: %w: %s", op, repo.ErrQueryExec, err)
		t.conf.log.Warn(retErr.Error())
		return retErr
	}

	return nil
}

const selectMembers = `
	SELECT 
		users.id, users.username, users.is_active
	FROM users 
		JOIN teams 
		ON users.team_id=teams.id
	WHERE name=$1
`

// getTeamMembers defines the logic of getting the members for the current team.
func (t teamsRepo) getTeamMembers(ctx context.Context, name string) ([]entities.User, error) {
	const op = "postgres.get-team-members"

	rows, err := t.conf.conn.Query(ctx, selectMembers, name)

	if err != nil {
		retErr := fmt.Errorf("error of the %s: %w: %w", op, repo.ErrQueryExec, err)
		t.conf.log.Warn(retErr.Error())
		return nil, retErr
	}
	defer rows.Close()

	res := make([]entities.User, 250)
	for rows.Next() {
		user := entities.User{}
		rows.Scan(&user.ID, &user.Name, &user.IsActive)
		res = append(res, user)
	}

	if rows.Err() != nil {
		retErr := fmt.Errorf("error of the %s: %w: %w", op, repo.ErrResProcessing, err)
		t.conf.log.Warn(retErr.Error())
		return nil, retErr
	}
	return res, nil
}

const selectTeam = `
	SELECT id, team_name
	FROM teams
	WHERE name=$1
`

// isTeamExists defines the logic of checking whether the current team exists.
func (t teamsRepo) isTeamExists(ctx context.Context, name string) (entities.Team, bool, error) {
	const op = "postgres.is-team-exists"

	rows, err := t.conf.conn.Query(ctx, selectTeam, name)

	if err != nil {
		retErr := fmt.Errorf("error of the %s: %w: %w", op, repo.ErrQueryExec, err)
		t.conf.log.Warn(retErr.Error())
		return entities.Team{}, false, retErr
	}
	defer rows.Close()

	res := entities.NewTeam()

	if rows.Next() {
		rows.Scan(&res.Name)
		return res, true, nil
	}

	if rows.Err() != nil {
		retErr := fmt.Errorf("error of the %s: %w: %w", op, repo.ErrResProcessing, rows.Err())
		t.conf.log.Warn(retErr.Error())
		return entities.Team{}, false, retErr
	}

	return entities.Team{}, false, nil
}
