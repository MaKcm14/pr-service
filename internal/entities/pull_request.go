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
	CreatedAt *time.Time
	MergedAt  *time.Time
	Author    User
	Reviewers map[UserID]User
}

func NewPullRequest() PullRequest {
	return PullRequest{
		Reviewers: make(map[UserID]User, 5),
	}
}

func (p *PullRequest) SetCreatedAtNow() {
	p.CreatedAt = new(time.Time)
	(*p.CreatedAt) = time.Now()
}

func (p *PullRequest) SetMergedAtNow() {
	p.MergedAt = new(time.Time)
	(*p.MergedAt) = time.Now()
}

func (p *PullRequest) SetReviewers(team Team) error {
	gen := makeReviewerRandGen(team.Members, nil)
	for i := 0; i != min(len(team.Members), rand.Intn(2)+1); {
		pos, err := gen()
		if err != nil {
			break
		}
		p.Reviewers[team.Members[pos].ID] = team.Members[pos]
		i++
	}

	if len(p.Reviewers) == 0 {
		return ErrReviewerAssign
	}
	return nil
}

func (p *PullRequest) ReassignReviewer(id UserID, team Team) (UserID, error) {
	if p.Status == Merged {
		return "", ErrStatusForReassign
	}

	if !p.CheckUserIsReviewer(id) {
		return "", ErrReviewerIsWrong
	}

	delete(p.Reviewers, id)

	gen := makeReviewerRandGen(team.Members, []UserID{id})

	pos, err := gen()
	if err != nil {
		return "", err
	}

	p.Reviewers[team.Members[pos].ID] = team.Members[pos]
	return team.Members[pos].ID, nil
}

func (p *PullRequest) CheckUserIsReviewer(id UserID) bool {
	_, val := p.Reviewers[id]
	return val
}
