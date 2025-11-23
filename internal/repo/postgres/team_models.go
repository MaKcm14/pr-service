package postgres

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/MaKcm14/pr-service/internal/entities"
	"github.com/MaKcm14/pr-service/internal/repo"
)

// teamsRepo defines the logic of interaction with the teams models.
type teamsRepo struct {
	conf *postgresConfig
}

// GetTeam defines the logic of getting the team object for the current name.
func (p *PostgreSQLRepo) GetTeam(ctx context.Context, name string) (entities.Team, error) {
	const op = "postgres.get-team"

	_, err := p.teamsRepo.isTeamExists(ctx, name)
	if err != nil {
		retErr := fmt.Errorf("error of the %s: %w", op, err)

		if errors.Is(err, repo.ErrModelNotFound) {
			return entities.Team{}, fmt.Errorf("error of the %s: %w", op, err)
		}
		p.conf.log.Warn(retErr.Error())

		return entities.Team{}, retErr
	}

	team := entities.NewTeam()
	team.Name = name

	members, err := p.teamsRepo.getTeamMembers(ctx, name)
	if err != nil {
		retErr := fmt.Errorf("error of the %s: %w", op, err)
		p.conf.log.Warn(retErr.Error())
		return entities.Team{}, retErr
	}
	team.Members = members

	return team, nil
}

// CreateTeam defines the logic of creating a new team.
func (p *PostgreSQLRepo) CreateTeam(ctx context.Context, team entities.Team) error {
	const op = "postgres.create-team"

	team, err := p.teamsRepo.isTeamExists(ctx, team.Name)
	if err != nil && !errors.Is(err, repo.ErrModelNotFound) {
		retErr := fmt.Errorf("error of the %s: %w", op, err)
		p.conf.log.Warn(retErr.Error())
		return retErr
	} else if err == nil {
		if err := p.teamsRepo.updateTeam(ctx, team); err != nil {
			retErr := fmt.Errorf("error of the %s: %w", op, err)
			p.conf.log.Warn(retErr.Error())
			return retErr
		}
		return fmt.Errorf("error of the %s: %w", op, repo.ErrModelAlreadyExists)
	}

	if err := p.teamsRepo.createTeam(ctx, team); err != nil {
		return err
	}

	return nil
}

func (t teamsRepo) updateTeam(ctx context.Context, team entities.Team) error {
	const op = "postgres.update-team"

	members, err := t.getTeamMembers(ctx, team.Name)
	if err != nil {
		retErr := fmt.Errorf("error of the %s: %w: %s", op, repo.ErrQueryExec, err)
		t.conf.log.Warn(retErr.Error())
		return retErr
	}

	updateList := make([]entities.User, 0, 100)
	addList := make([]entities.User, 0, 100)

	for _, newUser := range team.Members {
		flagExists := false

		for _, user := range members {
			if user.ID == newUser.ID {
				flagExists = true
				updateList = append(updateList, newUser)
				break
			}
		}

		if !flagExists {
			addList = append(addList, newUser)
		}
	}

	if err := t.addMembersList(ctx, addList, team); err != nil {
		retErr := fmt.Errorf("error of the %s: %w", op, err)
		return retErr
	}

	if err := t.updateMembers(ctx, updateList, team.ID); err != nil {
		retErr := fmt.Errorf("error of the %s: %w", op, err)
		return retErr
	}

	return nil
}

const updateMembers = `
	UPDATE users
	SET username=$1, is_active=$2, team_id=$3
	WHERE id = $4
`

func (t teamsRepo) updateMembers(
	ctx context.Context,
	members []entities.User,
	teamID entities.TeamID,
) error {
	const op = "postgres.update-members"

	for _, user := range members {
		_, err := t.conf.conn.Exec(ctx, updateMembers, user.Name,
			getSqlViewBool(user.IsActive), teamID, user.ID)

		if err != nil {
			retErr := fmt.Errorf("error of the %s: %w: %s", op, repo.ErrQueryExec, err)
			t.conf.log.Warn(retErr.Error())
			return retErr
		}
	}

	return nil
}

const insertTeam = `
	INSERT INTO teams (team_name)
	VALUES ($1)
	RETURNING id
`

// createTeam defines the logic of creating the 'Team' model.
func (t teamsRepo) createTeam(ctx context.Context, team entities.Team) error {
	const op = "postgres.create-team"

	rows, err := t.conf.conn.Query(ctx, insertTeam, team.Name)

	if err != nil {
		retErr := fmt.Errorf("error of the %s: %w: %w", op, repo.ErrQueryExec, err)
		t.conf.log.Warn(retErr.Error())
		return retErr
	}
	defer rows.Close()

	if rows.Next() {
		rows.Scan(&team.ID)
	} else if rows.Err() != nil {
		retErr := fmt.Errorf("error of the %s: %w: %w", op, repo.ErrResProcessing, err)
		t.conf.log.Warn(retErr.Error())
		return retErr
	}

	if err := t.addMembersList(ctx, team.Members, team); err != nil {
		return err
	}
	return nil
}

const insertMembers = `
INSERT INTO users (id, username, is_active, team_id)
VALUES `

// addMembersList defines the logic of adding the members list for the current team.
func (t teamsRepo) addMembersList(ctx context.Context, list []entities.User, team entities.Team) error {
	const op = "postgres.add-members-list"

	if len(team.Members) == 0 {
		return nil
	}

	query := bytes.Buffer{}
	query.WriteString(insertMembers)
	for idx, user := range list {
		query.WriteString(fmt.Sprintf("('%s', '%s', %s, %d)", string(user.ID), user.Name,
			getSqlViewBool(user.IsActive), team.ID))

		if idx != len(list)-1 {
			query.WriteString(",\n")
		}
	}

	_, err := t.conf.conn.Exec(ctx, query.String())

	fmt.Println(query.String())

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
	WHERE teams.team_name=$1
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

	res := make([]entities.User, 0, 250)
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
	WHERE team_name=$1
`

// isTeamExists defines the logic of checking whether the current team exists.
func (t teamsRepo) isTeamExists(ctx context.Context, name string) (entities.Team, error) {
	const op = "postgres.is-team-exists"

	res := entities.Team{}
	rows, err := t.conf.conn.Query(ctx, selectTeam, name)
	if err != nil {
		retErr := fmt.Errorf("error of the %s: %w: %w", op, repo.ErrQueryExec, err)
		t.conf.log.Warn(retErr.Error())
		return entities.Team{}, retErr
	}
	defer rows.Close()

	if rows.CommandTag().RowsAffected() != 0 {
		if rows.Next() {
			rows.Scan(&res.ID, &res.Name)
		}

		if rows.Err() != nil {
			retErr := fmt.Errorf("error of the %s: %w: %w", op, repo.ErrResProcessing, err)
			t.conf.log.Warn(retErr.Error())
			return entities.Team{}, retErr
		}
		return res, nil
	}

	return entities.Team{}, fmt.Errorf("error of the %s: %w", op, repo.ErrModelNotFound)
}
