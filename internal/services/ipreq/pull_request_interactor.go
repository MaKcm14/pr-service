package ipreq

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/MaKcm14/pr-service/internal/entities/dto"
	"github.com/MaKcm14/pr-service/internal/repo"
	"github.com/MaKcm14/pr-service/internal/services"
)

// PullRequestInteractor defines the logic of the use-cases more connected with the pull-requests.
type PullRequestInteractor struct {
	log  *slog.Logger
	repo services.PullRequestRepository
}

func NewPullRequestInteractor(log *slog.Logger, repo services.PullRequestRepository) PullRequestInteractor {
	return PullRequestInteractor{
		log:  log,
		repo: repo,
	}
}

// CreatePullRequest defines the logic of creating the pull-request.
func (p PullRequestInteractor) CreatePullRequest(ctx context.Context, pullRequest dto.PullRequestDTO) error {
	const op = "ipreq.create-pull-request"

	if err := p.repo.CreatePullRequest(ctx, pullRequest); err != nil {
		retErr := fmt.Errorf("error of the %s: %w: %s", op, services.ErrRepositoryInteraction, err)

		if errors.Is(err, repo.ErrModelNotFound) {
			return fmt.Errorf("error of the %s: %w: %s", op, services.ErrEntityNotFound, err)
		} else if errors.Is(err, repo.ErrModelAlreadyExists) {
			return fmt.Errorf("error of the %s: %w: %s", op, services.ErrEntityAlreadyExists, err)
		}
		p.log.Warn(retErr.Error())

		return retErr
	}

	return nil
}
