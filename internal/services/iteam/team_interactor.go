package iteam

import (
	"fmt"
	"log/slog"

	"github.com/MaKcm14/pr-service/internal/entities"
	"github.com/MaKcm14/pr-service/internal/services"
)

// TeamUseCase defines the logic of the use-cases more connected with the teams.
type TeamUseCase struct {
	log  *slog.Logger
	repo services.Repository
}

func NewTeamUseCase(log *slog.Logger, repo services.Repository) TeamUseCase {
	return TeamUseCase{
		log:  log,
		repo: repo,
	}
}

// GetTeam defines the logic of getting the team from the repository.
func (t TeamUseCase) GetTeam(name string) (entities.Team, bool, error) {
	const op = "iteam.get-team"

	dto, isExists, err := t.repo.GetTeam(name)
	if err != nil {
		retErr := fmt.Errorf("error of the %s: %w", op, err)
		t.log.Warn(retErr.Error())
		return entities.Team{}, false, retErr
	}

	return dto, isExists, nil
}

// CreateTeam defines the logic of creating the team in the repository.
func (t TeamUseCase) CreateTeam(team entities.Team) error {
	const op = "iteam.create-team"

	if err := t.repo.CreateTeam(team); err != nil {
		retErr := fmt.Errorf("error of the %s: %w", op, err)
		t.log.Warn(retErr.Error())
		return retErr
	}

	return nil
}
