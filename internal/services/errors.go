package services

import "errors"

var (
	ErrRepositoryInteraction  = errors.New("services: error of the interaction with the repository")
	ErrEntityNotFound         = errors.New("services: error of entity search: it wasn't found")
	ErrEntityAlreadyExists    = errors.New("services: entitiy already exists")
	ErrDomainRulesWithROState = errors.New("services: error of the domain's rules: can't complete the current operation due to its RO-state for this entity")
	ErrDomainRulesNoCandidate = errors.New("services: error of the finding the needed candidates")
	ErrWrongCandidate         = errors.New("services: error of using the current candidate")
)
