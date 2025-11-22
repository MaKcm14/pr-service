package entities

import "time"

const (
	Open   PullRequestStatus = "OPEN"
	Merged PullRequestStatus = "MERGED"
)

type PullRequestID string

// PullRequestStatus defines the common pull-request status.
type PullRequestStatus string

// PullRequest defines the entities object describes the pull-requests.
type PullRequest struct {
	ID        PullRequestID
	Name      string
	Status    PullRequestStatus
	CreatedAt time.Time
	MergedAt  time.Time
	Author    User
	Reviewers []User
}

func NewPullRequest() PullRequest {
	return PullRequest{
		Reviewers: make([]User, 0, 5),
	}
}

func (p *PullRequest) SetCreatedAtNow() {
	p.CreatedAt = time.Now()
}

func (p *PullRequest) SetMergedAtNow() {
	p.MergedAt = time.Now()
}
