package chttp

import "errors"

var (
	ErrQueryParam      = errors.New("chttp: error of the query param's domain")
	ErrQueryEmptyParam = errors.New("chttp: error while processing the empty param")
)

var (
	ErrRespQueryEmptyParam      = errors.New("the requested params couldn't be empty")
	ErrRespQueryUnexistingModel = errors.New("the requested model doesn't exist")
)
