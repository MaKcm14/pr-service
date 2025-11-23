package repo

import "errors"

var (
	ErrConnToRepository           = errors.New("repo: error of connection to the repository")
	ErrQueryExec                  = errors.New("repo: error of execution the query")
	ErrResProcessing              = errors.New("repo: error of processing the result")
	ErrCreateMultipleUniqueModels = errors.New("repo: error of creating the multiple unique models")
	ErrModelNotFound              = errors.New("repo: error of finding the model")
	ErrDependModelsNotFound       = errors.New("repo: error of finding the dependendent model")
	ErrModelAlreadyExists         = errors.New("repo: model already exists")
	ErrStartTransaction           = errors.New("repo: error of starting the transaction")
)
