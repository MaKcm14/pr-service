package chttp

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/MaKcm14/pr-service/internal/entities"
	"github.com/MaKcm14/pr-service/internal/repo"
	"github.com/MaKcm14/pr-service/internal/services"
	"github.com/labstack/echo/v4"
)

// HttpController defines the logic defining the requests handling process.
type HttpController struct {
	log     *slog.Logger
	server  *echo.Echo
	useCase services.Interactor
}

func New(log *slog.Logger, socket string, interactor services.Interactor) HttpController {
	contr := HttpController{
		log:     log,
		server:  echo.New(),
		useCase: interactor,
	}
	contr.configEndpoints()

	return contr
}

// configEndpoints sets the endpoints for the current kernel server instance.
func (h *HttpController) configEndpoints() {
	h.server.GET("/team/get", h.handlerTeamAdd)
	h.server.GET("/users/getReview", h.handlerUsersGetReview)

	h.server.POST("/team/add", h.handlerTeamAdd)
	h.server.POST("/users/setIsActive", h.handlerUserSetIsActive)
	h.server.POST("/pullRequest/create", h.handlerPullRequestCreate)
	h.server.POST("/pullRequest/merge", h.handlerPullRequestMerge)
	h.server.POST("/pullRequest/reassign", h.handlerPullRequestReassign)
}

// handlerTeamGet defines the logic of handling the request for getting the team.
func (h *HttpController) handlerTeamGet(eCtx echo.Context) error {
	const op = "chttp.team-get"

	res, err := validateTeamName(eCtx)
	if err != nil {
		return eCtx.JSON(http.StatusBadRequest, ErrResponse{
			ErrData{
				Code:    NoCandidate,
				Message: ErrRespQueryEmptyParam.Error(),
			}})
	}

	ctx, _ := context.WithTimeout(eCtx.Request().Context(), time.Second*3)

	dto, isExists, err := h.useCase.GetTeam(ctx, res.(string))
	if err != nil {
		h.log.Warn("error of the %s: %s", op, err)
		return eCtx.JSON(http.StatusOK, dto)
	} else if !isExists {
		return eCtx.JSON(http.StatusNotFound, ErrResponse{
			ErrData{
				Code:    NotFound,
				Message: ErrRespQueryNotFound.Error(),
			}})
	}

	return nil
}

func (h *HttpController) handlerUsersGetReview(ctx echo.Context) error {
	const op = "chttp.users-get-review"
	return nil
}

// handlerTeamAdd defines the logic of handling the request for adding the team.
func (h *HttpController) handlerTeamAdd(eCtx echo.Context) error {
	const op = "chttp.team-add"

	team := entities.NewTeam()

	if err := eCtx.Bind(&team); err != nil {
		return eCtx.JSON(http.StatusBadRequest, ErrResponse{
			ErrData{
				Code:    NoCandidate,
				Message: ErrRespQueryWrongRequestData.Error(),
			}})
	}

	ctx, _ := context.WithTimeout(eCtx.Request().Context(), time.Second*3)

	if err := h.useCase.CreateTeam(ctx, team); err != nil {
		if errors.Is(err, repo.ErrCreateMultipleUniqueModels) {
			return eCtx.JSON(http.StatusBadRequest, ErrResponse{
				ErrData{
					Code:    TeamExists,
					Message: ErrRespQueryAlreadyExists.Error(),
				}})
		}
	}

	return nil
}

func (h *HttpController) handlerUserSetIsActive(ctx echo.Context) error {
	const op = "chttp.user-set-is-active"
	return nil
}

func (h *HttpController) handlerPullRequestCreate(ctx echo.Context) error {
	const op = "chttp.pull-request-create"
	return nil
}

func (h *HttpController) handlerPullRequestMerge(ctx echo.Context) error {
	const op = "chttp.pull-request-merge"
	return nil
}

func (h *HttpController) handlerPullRequestReassign(ctx echo.Context) error {
	const op = "chttp.pull-request-reassign"
	return nil
}
