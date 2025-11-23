package services

import (
	"context"

	"github.com/MaKcm14/pr-service/internal/entities"
	"github.com/MaKcm14/pr-service/internal/entities/dto"
)

type (
	// TeamRepository defines the abstraction of the team's model ops interaction.
	TeamRepository interface {
		Closer

		GetTeam(ctx context.Context, name string) (entities.Team, error)
		CreateTeam(ctx context.Context, team entities.Team) error
	}

	// UserRepository defines the abstraction of the user's model ops interaction.
	UserRepository interface {
		Closer

		SetUserIsActive(ctx context.Context, isActive bool, id entities.UserID) (entities.User, error)
		GetUser(ctx context.Context, id entities.UserID) (entities.User, error)
	}

	PullRequestRepository interface {
		Closer

		CreatePullRequest(ctx context.Context, pullRequest dto.PullRequestDTO) error
		SetPullRequestStatus(ctx context.Context, status entities.PullRequestStatus, pullReq dto.PullRequestDTO) (dto.PullRequestDTO, error)
		GetUserPullRequests(ctx context.Context, id entities.UserID) ([]dto.PullRequestDTOShort, error)
		GetPullRequest(ctx context.Context, id entities.PullRequestID) (dto.PullRequestDTO, error)
		ChangeReviewer(ctx context.Context, lastID entities.UserID, newID entities.UserID, pullReq dto.PullRequestDTO) error
	}

	Closer interface {
		Close()
	}
)
