package ipreq

import "log/slog"

// PullRequestInteractor defines the logic of the use-cases more connected with the pull-requests.
type PullRequestInteractor struct {
	log *slog.Logger
}

func NewPullRequestInteractor(log *slog.Logger) PullRequestInteractor {
	return PullRequestInteractor{
		log: log,
	}
}
