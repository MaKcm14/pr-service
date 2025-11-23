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
		Closer

		GetTeam(ctx context.Context, teamName string) (dto.TeamDTO, error)
		CreateTeam(ctx context.Context, dto entities.Team) error
	}

	// UserInteractor defines the interface of the user's use-cases abstraction.
	UserInteractor interface {
		Closer

		SetUserIsActive(ctx context.Context, isActive bool, id entities.UserID) (entities.User, error)
	}

	// PullRequestInteractor defines the interface of the pull-requests' user-cases abstraction.
	PullRequestInteractor interface {
		Closer

		CreatePullRequest(ctx context.Context, pullReq dto.PullRequestDTO) error
		SetPullRequestStatus(ctx context.Context, status entities.PullRequestStatus, pullReq dto.PullRequestDTO) (dto.PullRequestDTO, error)
		GetUserPullRequests(ctx context.Context, id entities.UserID) ([]dto.PullRequestDTOShort, error)
		ReassignUser(ctx context.Context, reassignData dto.PullRequestChangeReviewerDTO) (dto.PullRequestDTO, entities.UserID, error)
	}
)
