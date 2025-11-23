package usecase

import (
	"log/slog"

	"github.com/MaKcm14/pr-service/internal/services"
	"github.com/MaKcm14/pr-service/internal/services/ipreq"
	"github.com/MaKcm14/pr-service/internal/services/iteam"
	"github.com/MaKcm14/pr-service/internal/services/iuser"
)

type UseCase struct {
	*ipreq.PullRequestUseCase
	*iteam.TeamUseCase
	*iuser.UserUseCase
}

func NewUseCase(
	log *slog.Logger,
	teamRepo services.TeamRepository,
	prRepo services.PullRequestRepository,
	userRepo services.UserRepository,
) UseCase {
	return UseCase{
		PullRequestUseCase: ipreq.NewPullRequestUseCase(log, prRepo, userRepo, teamRepo),
		TeamUseCase:        iteam.NewTeamUseCase(log, teamRepo),
		UserUseCase:        iuser.NewUserUseCase(log, userRepo),
	}
}

func (u UseCase) Close() {
	u.PullRequestUseCase.Close()
	u.TeamUseCase.Close()
	u.UserUseCase.Close()
}
