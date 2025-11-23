package entities

import (
	"math/rand"
	"time"
)

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

func (p *PullRequest) SetReviewers(team Team) {
	buff := make(map[int]struct{})

	for i := 0; i != min(len(team.Members), rand.Intn(2)+1); {
		idx := rand.Intn(len(team.Members))

		if _, ok := buff[idx]; !ok {
			if team.Members[idx].IsActive {
				p.Reviewers = append(p.Reviewers, team.Members[idx])
			}
			buff[idx] = struct{}{}
			i++
		}
	}
}
