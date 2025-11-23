package ipreq

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/MaKcm14/pr-service/internal/entities"
	"github.com/MaKcm14/pr-service/internal/entities/dto"
	"github.com/MaKcm14/pr-service/internal/repo"
	"github.com/MaKcm14/pr-service/internal/services"
)

// PullRequestInteractor defines the logic of the use-cases more connected with the pull-requests.
type PullRequestInteractor struct {
	log      *slog.Logger
	prRepo   services.PullRequestRepository
	userRepo services.UserRepository
	teamRepo services.TeamRepository
}

func NewPullRequestInteractor(
	log *slog.Logger,
	prRepo services.PullRequestRepository,
	userRepo services.UserRepository,
	teamRepo services.TeamRepository,
) PullRequestInteractor {
	return PullRequestInteractor{
		log:      log,
		prRepo:   prRepo,
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}

// CreatePullRequest defines the logic of creating the pull-request.
func (p PullRequestInteractor) CreatePullRequest(ctx context.Context, pullRequest dto.PullRequestDTO) error {
	const op = "ipreq.create-pull-request"

	user, err := p.userRepo.GetUser(ctx, pullRequest.AuthorID)
	if err != nil {
		retErr := fmt.Errorf("error of the %s: %w: %s", op, services.ErrRepositoryInteraction, err)

		if errors.Is(err, repo.ErrModelNotFound) {
			return fmt.Errorf("error of the %s: %w: %s", op, services.ErrEntityNotFound, err)
		}
		p.log.Warn(retErr.Error())

		return retErr
	}

	team, err := p.teamRepo.GetTeam(ctx, user.TeamName)
	if err != nil {
		retErr := fmt.Errorf("error of the %s: %w: %s", op, services.ErrRepositoryInteraction, err)
		p.log.Warn(retErr.Error())
		return retErr
	}
	pullReq := dto.PullRequestDTOToPullRequest(pullRequest)
	pullReq.SetReviewers(team)

	if err := p.prRepo.CreatePullRequest(ctx, pullRequest); err != nil {
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

// SetPullRequestStatus defines the logic of changing the status for the pull-request object.
func (p PullRequestInteractor) SetPullRequestStatus(
	ctx context.Context,
	status entities.PullRequestStatus,
	pullReq dto.PullRequestDTO,
) (dto.PullRequestDTO, error) {
	const op = "ipreq.set-pull-request-status"

	res, err := p.prRepo.SetPullRequestStatus(ctx, status, pullReq)
	if err != nil {
		retErr := fmt.Errorf("error of the %s: %w: %s", op, services.ErrRepositoryInteraction, err)

		if errors.Is(err, repo.ErrModelNotFound) {
			return dto.PullRequestDTO{},
				fmt.Errorf("error of the %s: %w: %s", op, services.ErrEntityNotFound, err)
		}
		p.log.Warn(retErr.Error())

		return dto.PullRequestDTO{}, retErr
	}

	return res, nil
}

func (p PullRequestInteractor) GetUserPullRequests(ctx context.Context, id entities.UserID) ([]dto.PullRequestDTOShort, error) {
	const op = "ipreq.get-user-pull-requests"

	res, err := p.prRepo.GetUserPullRequests(ctx, id)

	if err != nil {
		retErr := fmt.Errorf("error of the %s: %w: %s", op, services.ErrRepositoryInteraction, err)

		if errors.Is(err, repo.ErrModelNotFound) {
			return nil, fmt.Errorf("error of the %s: %w: %s", op, services.ErrEntityNotFound, err)
		}
		p.log.Warn(retErr.Error())

		return nil, retErr
	}

	return res, nil
}

func (p PullRequestInteractor) ReassignUser(ctx context.Context, reassignData dto.PullRequestChangeReviewerDTO) (dto.PullRequestDTO, entities.UserID, error) {
	const op = "ipreq.reassign-user"

	user, err := p.userRepo.GetUser(ctx, reassignData.OldUserID)
	if err != nil {
		retErr := fmt.Errorf("error of the %s: %w: %s", op, services.ErrRepositoryInteraction, err)

		if errors.Is(err, repo.ErrModelNotFound) {
			return dto.PullRequestDTO{}, "", fmt.Errorf("error of the %s: %w: %s", op, services.ErrEntityNotFound, err)
		}
		p.log.Warn(retErr.Error())

		return dto.PullRequestDTO{}, "", retErr
	}

	pullReq, err := p.prRepo.GetPullRequest(ctx, reassignData.ID)
	if err != nil {
		retErr := fmt.Errorf("error of the %s: %w: %s", op, services.ErrRepositoryInteraction, err)

		if errors.Is(err, repo.ErrModelNotFound) {
			return dto.PullRequestDTO{}, "", fmt.Errorf("error of the %s: %w: %s", op, services.ErrEntityNotFound, err)
		}
		p.log.Warn(retErr.Error())

		return dto.PullRequestDTO{}, "", retErr
	}

	team, err := p.teamRepo.GetTeam(ctx, user.TeamName)
	if err != nil {
		retErr := fmt.Errorf("error of the %s: %w: %s", op, services.ErrRepositoryInteraction, err)

		if errors.Is(err, repo.ErrModelNotFound) {
			return dto.PullRequestDTO{}, "", fmt.Errorf("error of the %s: %w: %s", op, services.ErrEntityNotFound, err)
		}
		p.log.Warn(retErr.Error())

		return dto.PullRequestDTO{}, "", retErr
	}

	prEnt := dto.PullRequestDTOToPullRequest(pullReq)
	id, err := prEnt.ReassignReviewer(user.ID, team)

	if err != nil {
		if errors.Is(err, entities.ErrStatusForReassign) {
			return dto.PullRequestDTO{}, "", fmt.Errorf("error of the %s: %w: %s", op, services.ErrDomainRulesWithROState, err)
		} else if errors.Is(err, entities.ErrReviewerAssign) {
			return dto.PullRequestDTO{}, "", fmt.Errorf("error of the %s: %w: %s", op, services.ErrDomainRulesNoCandidate, err)
		} else if errors.Is(err, entities.ErrReviewerIsWrong) {
			return dto.PullRequestDTO{}, "", fmt.Errorf("error of the %s: %w: %s", op, services.ErrWrongCandidate, err)
		}
	}

	if err := p.prRepo.ChangeReviewer(ctx, reassignData.OldUserID, id, pullReq); err != nil {
		retErr := fmt.Errorf("error of the %s: %w: %s", op, services.ErrRepositoryInteraction, err)
		p.log.Warn(retErr.Error())
		return dto.PullRequestDTO{}, "", retErr
	}

	return dto.PullRequestToPullRequestDTO(prEnt), id, nil
}
