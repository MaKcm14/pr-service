package chttp

import "errors"

var (
	ErrStartingServer = errors.New("chttp: error of starting the server")
)

var (
	ErrQueryParam      = errors.New("chttp: error of the query param's domain")
	ErrQueryEmptyParam = errors.New("chttp: error while processing the empty param")
)

var (
	ErrRespQueryEmptyParam       = errors.New("the requested params couldn't be empty")
	ErrRespQueryNotFound         = errors.New("the requested model doesn't exist")
	ErrRespQueryWrongRequestData = errors.New("error of the request's data: current format is not available")
	ErrRespQueryAlreadyExists    = errors.New("the model already exists")
	ErrRespQueryServerError      = errors.New("internal error was generated during request processing")
	ErrRespQueryOpIsRestrict     = errors.New("can't complete this operation due to the violation of changing the internal state")
	ErrRespQueryNoCandidate      = errors.New("couldn't find the needed candidates for the current operation")
	ErrRespQueryWrongCandidate   = errors.New("couldn't complete the operation with the current candidate due to it's wrong")
)
