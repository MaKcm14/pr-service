package chttp

import (
	"log/slog"
	"net/http"

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
func (h *HttpController) handlerTeamGet(ctx echo.Context) error {
	const op = "chttp.team-get"

	res, err := validateTeamName(ctx)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, ErrResponse{
			ErrRespQueryEmptyParam.Error(),
		})
	}

	team, isExists, err := h.useCase.GetTeam(res.(string))
	if err != nil {
		h.log.Warn("error of the %s: %s", op, err)
		return ctx.JSON(http.StatusOK, team)
	} else if !isExists {
		return ctx.JSON(http.StatusNotFound, ErrResponse{
			ErrRespQueryUnexistingModel.Error(),
		})
	}

	return nil
}

func (h *HttpController) handlerUsersGetReview(ctx echo.Context) error {
	const op = "chttp.users-get-review"
	return nil
}

func (h *HttpController) handlerTeamAdd(ctx echo.Context) error {
	const op = "chttp.team-add"

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
