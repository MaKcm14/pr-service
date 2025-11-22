package chttp

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

// ValidFunc defines the validating-func template.
type ValidFunc func(eCtx echo.Context) (any, error)

// validateTeamName checks whether the current team-name is correct.
func validateTeamName(eCtx echo.Context) (any, error) {
	const op = "chttp.validate-team-name"

	name := eCtx.QueryParam("team_name")
	if len(name) == 0 {
		return nil, fmt.Errorf("error of the %s: %w: %w", ErrQueryParam, ErrQueryEmptyParam)
	}
	return name, nil
}
