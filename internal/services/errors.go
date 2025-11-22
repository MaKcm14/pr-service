package services

import "errors"

var (
	ErrRepositoryInteraction = errors.New("services: error of the interaction with the repository")
	ErrEntityNotFound        = errors.New("services: error of entity search: it wasn't found")
	ErrEntityAlreadyExists   = errors.New("services: entitiy already exists")
)
