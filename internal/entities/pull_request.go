package entities

import "time"

const (
	Open   PullRequestStatus = "OPEN"
	Merged PullRequestStatus = "MERGED"
)

// PullRequestStatus defines the common pull-request status.
type PullRequestStatus string

// PullRequest defines the entities object describes the pull-requests.
type PullRequest struct {
	Name      string
	Status    PullRequestStatus
	CreatedAt time.Time
	MergedAt  time.Time
	Author    User
}

func (p *PullRequest) SetCreatedAtNow() {
	p.CreatedAt = time.Now()
}

func (p *PullRequest) SetMergedAtNow() {
	p.MergedAt = time.Now()
}
