package dto

import (
	"time"

	"github.com/MaKcm14/pr-service/internal/entities"
)

type PullRequestDTO struct {
	ID        entities.PullRequestID     `json:"pull_request_id"`
	Name      string                     `json:"pull_request_name"`
	Status    entities.PullRequestStatus `json:"status"`
	CreatedAt time.Time                  `json:"created_at"`
	MergedAt  time.Time                  `json:"merged_at"`
	AuthorID  entities.UserID            `json:"author_id"`
	Reviewers []entities.User            `json:"assigned_reviewers"`
}

func PullRequestToPullRequestDTO(pullRequest entities.PullRequest) PullRequestDTO {
	return PullRequestDTO{
		ID:        pullRequest.ID,
		Name:      pullRequest.Name,
		Status:    pullRequest.Status,
		CreatedAt: pullRequest.CreatedAt,
		MergedAt:  pullRequest.MergedAt,
		AuthorID:  pullRequest.Author.ID,
		Reviewers: append(make([]entities.User, 0), pullRequest.Reviewers...),
	}
}

type PullRequestDTOShort struct {
	ID       entities.PullRequestID     `json:"pull_request_id"`
	Name     string                     `json:"pull_request_name"`
	Status   entities.PullRequestStatus `json:"status"`
	AuthorID entities.UserID            `json:"author_id"`
}

func PullRequestToPullRequestDTOShort(pullRequest entities.PullRequest) PullRequestDTOShort {
	return PullRequestDTOShort{
		ID:       pullRequest.ID,
		Name:     pullRequest.Name,
		Status:   pullRequest.Status,
		AuthorID: pullRequest.Author.ID,
	}
}

func MakePullRequestDTOShort(pullRequest PullRequestDTO) PullRequestDTOShort {
	return PullRequestDTOShort{
		ID:       pullRequest.ID,
		Name:     pullRequest.Name,
		Status:   pullRequest.Status,
		AuthorID: pullRequest.AuthorID,
	}
}
