package services

import (
	"context"

	"github.com/MaKcm14/pr-service/internal/entities"
	"github.com/MaKcm14/pr-service/internal/entities/dto"
)

type (
	// Interactor defines the common interface for every interactor's abstraction.
	Interactor interface {
		TeamInteractor
		UserInteractor
		PullRequestInteractor
	}

	// TeamInteractor defines the interface of the teams's use-cases abstraction.
	TeamInteractor interface {
		GetTeam(ctx context.Context, teamName string) (dto.TeamDTO, error)
		CreateTeam(ctx context.Context, dto entities.Team) error
	}

	// UserInteractor defines the interface of the user's use-cases abstraction.
	UserInteractor interface {
		SetUserIsActive(ctx context.Context, dto entities.User) (entities.User, error)
	}

	// PullRequestInteractor defines the interface of the pull-requests' user-cases abstraction.
	PullRequestInteractor interface {
		CreatePullRequest(ctx context.Context, pullReq dto.PullRequestDTO) error
		SetPullRequestStatus(ctx context.Context, status entities.PullRequestStatus, pullReq dto.PullRequestDTO) (dto.PullRequestDTO, error)
		GetUserPullRequests(ctx context.Context, id entities.UserID) ([]dto.PullRequestDTOShort, error)
	}
)
