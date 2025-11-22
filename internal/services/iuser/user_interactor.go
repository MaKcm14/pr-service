package iuser

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/MaKcm14/pr-service/internal/entities"
	"github.com/MaKcm14/pr-service/internal/repo"
	"github.com/MaKcm14/pr-service/internal/services"
)

type UserUseCase struct {
	log  *slog.Logger
	repo services.UserRepository
}

func NewUserUseCase(log *slog.Logger, repo services.UserRepository) UserUseCase {
	return UserUseCase{
		log:  log,
		repo: repo,
	}
}

func (u UserUseCase) SetUserIsActive(ctx context.Context, dto entities.User) (entities.User, error) {
	const op = "iuser.set-user-is-active"

	user, err := u.repo.SetUserIsActive(ctx, dto)
	if err != nil {
		if errors.Is(err, repo.ErrModelNotFound) {
			return entities.User{}, fmt.Errorf(
				"error of the %s: %w: %w", op, services.ErrEntityNotFound, err,
			)
		}
		retErr := fmt.Errorf("error of the %s: %w: %w", op, err, services.ErrRepositoryInteraction)
		u.log.Warn(retErr.Error())
		return entities.User{}, retErr
	}

	return user, nil
}
