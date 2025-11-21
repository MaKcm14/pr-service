package postgres

import "errors"

var (
	ErrConnToRepository = errors.New("postgres: error of connection to the repository")
	ErrQueryExec        = errors.New("postgres: error of execution the query")
	ErrResProcessing    = errors.New("postgres: error of processing the result")
)
