package iteam

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

// TeamUseCase defines the logic of the use-cases more connected with the teams.
type TeamUseCase struct {
	log  *slog.Logger
	repo services.TeamRepository
}

func NewTeamUseCase(log *slog.Logger, repo services.TeamRepository) *TeamUseCase {
	return &TeamUseCase{
		log:  log,
		repo: repo,
	}
}

// GetTeam defines the logic of getting the team from the repository.
func (t *TeamUseCase) GetTeam(ctx context.Context, teamName string) (dto.TeamDTO, error) {
	const op = "iteam.get-team"

	team, err := t.repo.GetTeam(ctx, teamName)
	if err != nil {
		retErr := fmt.Errorf("error of the %s: %w: %s", op, services.ErrRepositoryInteraction, err)

		if errors.Is(err, repo.ErrModelNotFound) {
			return dto.TeamDTO{}, fmt.Errorf("error of the %s: %w: %s", op, services.ErrEntityNotFound, err)
		}
		t.log.Warn(retErr.Error())

		return dto.TeamDTO{}, retErr
	}

	return dto.TeamToTeamDTO(team), nil
}

// CreateTeam defines the logic of creating the team in the repository.
func (t *TeamUseCase) CreateTeam(ctx context.Context, dto entities.Team) error {
	const op = "iteam.create-team"

	if err := t.repo.CreateTeam(ctx, dto); err != nil {
		retErr := fmt.Errorf("error of the %s: %w", op, err)
		t.log.Warn(retErr.Error())
		return retErr
	}

	return nil
}

func (t *TeamUseCase) Close() {
	t.repo.Close()
}
