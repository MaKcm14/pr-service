package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/MaKcm14/pr-service/internal/entities"
	"github.com/MaKcm14/pr-service/internal/entities/dto"
	"github.com/MaKcm14/pr-service/internal/repo"
)

// pullRequestRepo defines the repo-object for interaction with the pull-requests models.
type pullRequestRepo struct {
	conf *postgresConfig
	usersRepo
}

// CreatePullRequest defines the logic of creating the pull request in the repo.
func (p pullRequestRepo) CreatePullRequest(ctx context.Context, pullRequest dto.PullRequestDTO) error {
	const op = "postgres.create-pull-request"

	if err := p.isExists(ctx, pullRequest.ID); err != nil && !errors.Is(err, repo.ErrModelNotFound) {
		retErr := fmt.Errorf("error of the %s: %w: %s", op, repo.ErrModelAlreadyExists, err)
		p.conf.log.Warn(retErr.Error())
		return retErr
	}

	if _, err := p.getUser(ctx, pullRequest.AuthorID); err != nil {
		retErr := fmt.Errorf("error of the %s: %w: %s", op, repo.ErrModelNotFound, err)
		return retErr
	}

	if err := p.createPullRequest(ctx, pullRequest); err != nil {
		retErr := fmt.Errorf("error of the %s: %w: %s", op, err)
		p.conf.log.Warn(retErr.Error())
		return retErr
	}
	return nil
}

const insertPullRequest = `
	INSERT INTO pull_requests (id, pr_name, status, created_at, merged_at, author_id)
	VALUES ($1, $2, $3, $4, $5, $6)
`

func (p pullRequestRepo) createPullRequest(
	ctx context.Context,
	pullRequest dto.PullRequestDTO,
) error {
	const op = "postgres.create-pull-request"

	_, err := p.conf.conn.Exec(
		ctx,
		insertPullRequest,
		pullRequest.ID,
		pullRequest.Name,
		pullRequest.Status,
		pullRequest.CreatedAt,
		pullRequest.MergedAt,
		pullRequest.AuthorID,
	)

	if err != nil {
		retErr := fmt.Errorf("error of the %s: %w: %s", op, repo.ErrQueryExec, err)
		p.conf.log.Warn(retErr.Error())
		return retErr
	}
	return nil
}

const checkExisting = `
	SELECT id
	FROM pull_requests
	WHERE id=$1
`

// isExists defines the logic of checking the existing of the pull request.
func (p pullRequestRepo) isExists(ctx context.Context, id entities.PullRequestID) error {
	const op = "postgres.is-exists"

	tag, err := p.conf.conn.Exec(ctx, checkExisting)
	if err != nil {
		retErr := fmt.Errorf("error of the %s: %w: %s", op, repo.ErrQueryExec, err)
		p.conf.log.Warn(retErr.Error())
		return retErr
	}

	if tag.RowsAffected() == 0 {
		return repo.ErrModelNotFound
	}
	return nil
}
