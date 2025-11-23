package services

import (
	"context"

	"github.com/MaKcm14/pr-service/internal/entities"
	"github.com/MaKcm14/pr-service/internal/entities/dto"
)

type (
	// Repository defines the common interface of the repo's model interaction.
	Repository interface {
		TeamRepository
		UserRepository
		PullRequestRepository
	}

	// TeamRepository defines the abstraction of the team's model ops interaction.
	TeamRepository interface {
		GetTeam(ctx context.Context, name string) (entities.Team, error)
		CreateTeam(ctx context.Context, team entities.Team) error
	}

	// UserRepository defines the abstraction of the user's model ops interaction.
	UserRepository interface {
		SetUserIsActive(ctx context.Context, dto entities.User) (entities.User, error)
		GetUser(ctx context.Context, id entities.UserID) (entities.User, error)
	}

	PullRequestRepository interface {
		CreatePullRequest(ctx context.Context, pullRequest dto.PullRequestDTO) error
		SetPullRequestStatus(ctx context.Context, status entities.PullRequestStatus, pullReq dto.PullRequestDTO) (dto.PullRequestDTO, error)
		GetUserPullRequests(ctx context.Context, id entities.UserID) ([]dto.PullRequestDTOShort, error)
	}
)
