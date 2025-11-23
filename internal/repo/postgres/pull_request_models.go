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
}

// CreatePullRequest defines the logic of creating the pull request in the repo.
func (p PostgreSQLRepo) CreatePullRequest(ctx context.Context, pullRequest dto.PullRequestDTO) error {
	const op = "postgres.create-pull-request"

	if err := p.isPullRequestExists(ctx, pullRequest.ID); err == nil || err != nil && !errors.Is(err, repo.ErrModelNotFound) {
		retErr := fmt.Errorf("error of the %s: %w", op, err)

		if err == nil {
			return fmt.Errorf("error of the %s: %w", op, repo.ErrModelAlreadyExists)
		}
		p.conf.log.Warn(retErr.Error())

		return retErr
	}

	_, err := p.GetUser(ctx, pullRequest.AuthorID)
	if err != nil {
		retErr := fmt.Errorf("error of the %s: %w", op, err)

		if errors.Is(err, repo.ErrModelNotFound) {
			return fmt.Errorf("error of the %s: %w", op, err)
		}
		p.conf.log.Warn(retErr.Error())

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
func (p pullRequestRepo) isPullRequestExists(ctx context.Context, id entities.PullRequestID) error {
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

const updatePRStatus = `
	UPDATE pull_requests
	SET status=$1
	WHERE id=$2
`

// SetPullRequestStatus defines the logic of changing the PR's status.
func (p pullRequestRepo) SetPullRequestStatus(
	ctx context.Context,
	status entities.PullRequestStatus,
	pullReq dto.PullRequestDTO,
) (dto.PullRequestDTO, error) {
	const op = "postgres.set-pull-request-to-merged"

	if err := p.isPullRequestExists(ctx, pullReq.ID); err != nil {
		retErr := fmt.Errorf("error of the %s: %w", op, err)

		if errors.Is(err, repo.ErrModelNotFound) {
			return dto.PullRequestDTO{}, retErr
		}
		p.conf.log.Warn(retErr.Error())

		return dto.PullRequestDTO{}, retErr
	}

	_, err := p.conf.conn.Exec(ctx, updatePRStatus, status, pullReq.Status)
	if err != nil {
		retErr := fmt.Errorf("error of the %s: %w: %s", op, repo.ErrQueryExec, err)
		p.conf.log.Warn(retErr.Error())
		return dto.PullRequestDTO{}, err
	}

	res, err := p.getPullRequest(ctx, pullReq.ID)
	if err != nil {
		retErr := fmt.Errorf("error of the %s: %w", op, err)
		p.conf.log.Warn(retErr.Error())
		return dto.PullRequestDTO{}, retErr
	}
	return res, nil
}

const selectPullRequest = `
	SELECT id, pr_name, status, created_at, merged_at, author_id
	FROM pull_requests
	WHERE id=$1
`

func (p PostgreSQLRepo) GetPullRequest(
	ctx context.Context,
	id entities.PullRequestID,
) (dto.PullRequestDTO, error) {
	const op = "postgres.get-pull-request-global"

	if err := p.isPullRequestExists(ctx, id); err != nil {
		retErr := fmt.Errorf("error of the %s: %w", op, err)

		if errors.Is(err, repo.ErrModelNotFound) {
			return dto.PullRequestDTO{}, retErr
		}
		p.conf.log.Warn(retErr.Error())

		return dto.PullRequestDTO{}, retErr
	}

	res, err := p.getPullRequest(ctx, id)
	if err != nil {
		retErr := fmt.Errorf("error of the %s: %w", op, err)
		p.conf.log.Warn(retErr.Error())
		return dto.PullRequestDTO{}, retErr
	}

	return res, nil
}

func (p pullRequestRepo) getPullRequest(
	ctx context.Context,
	id entities.PullRequestID,
) (dto.PullRequestDTO, error) {
	const op = "postgres.get-pull-request"

	rows, err := p.conf.conn.Query(ctx, selectMembers, id)
	if err != nil {
		retErr := fmt.Errorf("error of the %s: %w: %s", op, repo.ErrQueryExec, err)
		p.conf.log.Warn(retErr.Error())
		return dto.PullRequestDTO{}, retErr
	}
	defer rows.Close()

	res := dto.NewPullRequestDTO()
	if rows.Next() {
		err := rows.Scan(&res.ID, &res.Name, &res.Status, &res.CreatedAt, &res.MergedAt, &res.AuthorID)
		if err != nil {
			retErr := fmt.Errorf("error of the %s: %w: %s", op, repo.ErrResProcessing, err)
			p.conf.log.Warn(retErr.Error())
			return dto.PullRequestDTO{}, retErr
		}
	} else if rows.Err() != nil {
		retErr := fmt.Errorf("error of the %s: %w: %s", op, repo.ErrResProcessing, rows.Err())
		p.conf.log.Warn(retErr.Error())
		return dto.PullRequestDTO{}, retErr
	}

	reviewers, err := p.getPullRequestReviewers(ctx, id)
	if err != nil {
		retErr := fmt.Errorf("error of the %s: %w", op, err)
		p.conf.log.Warn(retErr.Error())
		return dto.PullRequestDTO{}, retErr
	}
	res.Reviewers = reviewers

	return res, nil
}

const selectPRReviewers = `
	SELECT assigned_reviewers.id
	FROM pull_requests 
		JOIN assigned_reviewers 
		ON pull_requests.id=assigned_reviewers.pr_id
	WHERE pull_requests.id=$1
`

func (p pullRequestRepo) getPullRequestReviewers(
	ctx context.Context,
	id entities.PullRequestID,
) ([]entities.UserID, error) {
	const op = "postgres.get-pull-request-reviewers"

	rows, err := p.conf.conn.Query(ctx, selectPRReviewers, id)
	if err != nil {
		retErr := fmt.Errorf("error of the %s: %w: %s", op, repo.ErrQueryExec, err)
		p.conf.log.Warn(retErr.Error())
		return nil, retErr
	}
	defer rows.Close()

	res := make([]entities.UserID, 0, 5)
	for rows.Next() {
		var id entities.UserID
		rows.Scan(&id)
		res = append(res, id)
	}

	if rows.Err() != nil {
		retErr := fmt.Errorf("error of the %s: %w: %s", op, repo.ErrResProcessing, err)
		p.conf.log.Warn(retErr.Error())
		return nil, retErr
	}

	return res, nil
}

const selectUserPRs = `
	SELECT pr.pr_id, pr.pr_name, pr.author_id, pr.status
	FROM pull_requests AS pr
	JOIN assigned_reviewers AS ar
	ON pr.pr_id=ar.pr_id
	WHERE ar.user_id=$1
`

func (p PostgreSQLRepo) GetUserPullRequests(ctx context.Context, id entities.UserID) ([]dto.PullRequestDTOShort, error) {
	const op = "postgres.get-user-pull-requests"

	if _, err := p.GetUser(ctx, id); err != nil {
		retErr := fmt.Errorf("error of the %s: %w", op, err)

		if errors.Is(err, repo.ErrModelNotFound) {
			return nil, retErr
		}
		p.conf.log.Warn(retErr.Error())

		return nil, retErr
	}

	rows, err := p.conf.conn.Query(ctx, selectUserPRs, id)
	if err != nil {
		retErr := fmt.Errorf("error of the %s: %w: %s", op, repo.ErrQueryExec, err)
		p.conf.log.Warn(retErr.Error())
		return nil, retErr
	}
	defer rows.Close()

	res := make([]dto.PullRequestDTOShort, 0, 20)
	for rows.Next() {
		pr := dto.PullRequestDTOShort{}

		if err := rows.Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status); err != nil {
			retErr := fmt.Errorf("error of the %s: %w: %s", op, repo.ErrResProcessing, err)
			p.conf.log.Warn(retErr.Error())
			return nil, retErr
		}
		res = append(res, pr)
	}

	if rows.Err() != nil {
		retErr := fmt.Errorf("error of the %s: %w: %s", op, repo.ErrResProcessing, err)
		p.conf.log.Warn(retErr.Error())
		return nil, retErr
	}

	return res, nil
}

const changeReviewer = `
	UPDATE assigned_reviewers
	SET user_id=$1
	WHERE pr_id=$1 AND user_id=$2
`

func (p PostgreSQLRepo) ChangeReviewer(ctx context.Context, lastID entities.UserID, newID entities.UserID, pullReq dto.PullRequestDTO) error {
	const op = "postgres.change-reviewers-list"

	tag, err := p.conf.conn.Exec(ctx, changeReviewer, lastID, pullReq.ID, newID)

	if err != nil {
		retErr := fmt.Errorf("error of the %s: %w: %s", op, repo.ErrQueryExec, err)
		p.conf.log.Warn(retErr.Error())
		return retErr
	}

	if tag.RowsAffected() == 0 {
		retErr := fmt.Errorf("error of the %s: %w: %s", op, repo.ErrModelNotFound, err)
		p.conf.log.Warn(retErr.Error())
		return retErr
	}

	return nil
}
