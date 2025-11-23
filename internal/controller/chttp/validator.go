package chttp

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

// ValidFunc defines the validating-func template.
type ValidFunc func(eCtx echo.Context) (any, error)

// validateTeamName checks whether the current team-name is correct.
func validateTeamName(eCtx echo.Context) (any, error) {
	name := eCtx.QueryParam("team_name")

	if len(name) == 0 {
		return nil, fmt.Errorf("error of the %s: %w", ErrQueryParam, ErrQueryEmptyParam)
	}

	return name, nil
}

func validateUserID(eCtx echo.Context) (any, error) {
	userID := eCtx.QueryParam("user_id")

	if len(userID) == 0 {
		return nil, fmt.Errorf("error of the %s: %w", ErrQueryParam, ErrQueryEmptyParam)
	}

	return userID, nil
}
