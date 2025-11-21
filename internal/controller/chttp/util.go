package chttp

const (
	TeamExists  ErrCode = "TEAM_EXISTS"
	PrExists    ErrCode = "PR_EXISTS"
	PrMerged    ErrCode = "PR_MERGED"
	NotAssigned ErrCode = "NOT_ASSIGNED"
	NoCandidate ErrCode = "NO_CANDIDATE"
	NotFound    ErrCode = "NOT_FOUND"
)

// ErrCode defines the string error's view description.
type ErrCode string

// ErrData defines the object describes the error's data.
type ErrData struct {
	Code    ErrCode `json:"code"`
	Message string  `json:"message"`
}

// ErrResponse defines the object that returns in case of errors.
type ErrResponse struct {
	Data ErrData `json:"error"`
}
