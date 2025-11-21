package services

import "github.com/MaKcm14/pr-service/internal/entities"

type (
	// Repository defines the common interface of the repo's interaction.
	Repository interface {
		TeamRepository
	}

	// TeamRepository defines the abstraction of the team's ops interaction.
	TeamRepository interface {
		GetTeam(name string) (entities.Team, bool, error)
		CreateTeam(team entities.Team) error
	}
)
