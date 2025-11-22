package chttp

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/MaKcm14/pr-service/internal/entities"
	"github.com/MaKcm14/pr-service/internal/entities/dto"
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
	h.server.GET("/team/get", h.handlerTeamGet)
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

	return eCtx.JSON(http.StatusOK, dto)
}

// handlerUsersGetReview defines the logic of handling the request for getting the reviewer's PRs.
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
		h.log.Warn(fmt.Sprintf("error of the %s: %s", op, err))
		return eCtx.JSON(http.StatusInternalServerError, ErrResponse{
			ErrData{
				Code:    ServerErr,
				Message: ErrRespQueryServerError.Error(),
			}})
	}

	return eCtx.JSON(http.StatusCreated, dto.TeamToTeamDTO(team))
}

// handlerUserSetIsActive defines the logic of handling the request for setting the user to the active.
func (h *HttpController) handlerUserSetIsActive(eCtx echo.Context) error {
	const op = "chttp.user-set-is-active"

	dto := entities.User{}
	if err := eCtx.Bind(&dto); err != nil {
		return eCtx.JSON(http.StatusBadRequest, ErrResponse{
			ErrData{
				Code:    NoCandidate,
				Message: ErrRespQueryWrongRequestData.Error(),
			}})
	}

	ctx, _ := context.WithTimeout(eCtx.Request().Context(), time.Second*3)

	user, err := h.useCase.SetUserIsActive(ctx, dto)
	if err != nil {
		if errors.Is(err, services.ErrEntityNotFound) {
			return eCtx.JSON(http.StatusNotFound, ErrResponse{
				ErrData{
					Code:    NotFound,
					Message: ErrRespQueryNotFound.Error(),
				}})
		}
		retErr := fmt.Errorf("error of the %s: %s", op, ErrRespQueryServerError, err)
		h.log.Warn(retErr.Error())
		return eCtx.JSON(http.StatusInternalServerError, ErrResponse{
			ErrData{
				Code:    ServerErr,
				Message: ErrRespQueryServerError.Error(),
			}})
	}

	return eCtx.JSON(http.StatusOK, user)
}

// handlerPullRequestCreate defines the logic of handling the request for creating the pull-request.
func (h *HttpController) handlerPullRequestCreate(eCtx echo.Context) error {
	const op = "chttp.pull-request-create"

	pullReq := dto.NewPullRequestDTO()
	if err := eCtx.Bind(&pullReq); err != nil {
		return eCtx.JSON(http.StatusBadRequest, ErrResponse{
			ErrData{
				Code:    RequestDataErr,
				Message: ErrRespQueryWrongRequestData.Error(),
			},
		})
	}

	ctx, _ := context.WithTimeout(eCtx.Request().Context(), time.Second*3)
	if err := h.useCase.CreatePullRequest(ctx, pullReq); err != nil {
		if errors.Is(err, services.ErrEntityNotFound) {
			return eCtx.JSON(http.StatusNotFound, ErrResponse{
				ErrData{
					Code:    NotFound,
					Message: ErrRespQueryNotFound.Error(),
				}})
		} else if errors.Is(err, services.ErrEntityAlreadyExists) {
			eCtx.JSON(http.StatusNotFound, ErrResponse{
				ErrData{
					Code:    PrExists,
					Message: ErrRespQueryAlreadyExists.Error(),
				}})
		}
		h.log.Warn(fmt.Sprintf("error of the %s: %s", op, err))
		return eCtx.JSON(http.StatusInternalServerError, ErrResponse{
			ErrData{
				Code:    ServerErr,
				Message: ErrRespQueryServerError.Error(),
			}})
	}

	return eCtx.JSON(http.StatusCreated, dto.MakePullRequestDTOShort(pullReq))
}

// handlerPullRequestMerge defines the logic of handling the request for merge the requested PR.
func (h *HttpController) handlerPullRequestMerge(eCtx echo.Context) error {
	const op = "chttp.pull-request-merge"

	pullReq := dto.NewPullRequestDTO()
	if err := eCtx.Bind(&pullReq); err != nil {
		return eCtx.JSON(http.StatusBadRequest, ErrResponse{
			ErrData{
				Code:    RequestDataErr,
				Message: ErrRespQueryWrongRequestData.Error(),
			},
		})
	}

	ctx, _ := context.WithTimeout(eCtx.Request().Context(), time.Second*3)
	res, err := h.useCase.SetPullRequestStatus(ctx, entities.Merged, pullReq)
	if err != nil {
		if errors.Is(err, services.ErrEntityNotFound) {
			return eCtx.JSON(http.StatusNotFound, ErrResponse{
				ErrData{
					Code:    NotFound,
					Message: ErrRespQueryNotFound.Error(),
				},
			})
		}
		h.log.Warn(fmt.Sprintf("error of the %s: %s", op, err))

		return eCtx.JSON(http.StatusInternalServerError, ErrResponse{
			ErrData{
				Code:    ServerErr,
				Message: ErrRespQueryServerError.Error(),
			},
		})
	}

	return eCtx.JSON(http.StatusOK, res)
}

func (h *HttpController) handlerPullRequestReassign(ctx echo.Context) error {
	const op = "chttp.pull-request-reassign"
	return nil
}
