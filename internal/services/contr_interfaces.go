package services

import "github.com/MaKcm14/pr-service/internal/entities"

type (
	// Interactor defines the common interface for every interactor's abstraction.
	Interactor interface {
		TeamInteractor
	}

	// TeamInteractor defines the interface of the teams's use-cases abstraction.
	TeamInteractor interface {
		GetTeam(name string) (entities.Team, bool, error)
		CreateTeam(team entities.Team) error
	}
)
