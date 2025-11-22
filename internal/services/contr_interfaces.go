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
	}

	// TeamInteractor defines the interface of the teams's use-cases abstraction.
	TeamInteractor interface {
		GetTeam(ctx context.Context, name string) (dto.TeamDTO, bool, error)
		CreateTeam(ctx context.Context, team entities.Team) error
	}
)
