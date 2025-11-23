package entities

import "errors"

var (
	ErrReviewerAssign    = errors.New("entities: error of assigning the reviewers")
	ErrStatusForReassign = errors.New("entities: error of reassigning with the current PR's status")
	ErrReviewerIsWrong   = errors.New("entities: the user is not in reviewer's list")
)
