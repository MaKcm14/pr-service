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
	Reviewers []entities.UserID          `json:"assigned_reviewers"`
}

func NewPullRequestDTO() PullRequestDTO {
	return PullRequestDTO{
		Reviewers: make([]entities.UserID, 0, 5),
	}
}

func PullRequestToPullRequestDTO(pullRequest entities.PullRequest) PullRequestDTO {
	dto := PullRequestDTO{
		ID:        pullRequest.ID,
		Name:      pullRequest.Name,
		Status:    pullRequest.Status,
		CreatedAt: pullRequest.CreatedAt,
		MergedAt:  pullRequest.MergedAt,
		AuthorID:  pullRequest.Author.ID,
		Reviewers: make([]entities.UserID, 0, len(pullRequest.Reviewers)),
	}

	for _, user := range pullRequest.Reviewers {
		dto.Reviewers = append(dto.Reviewers, user.ID)
	}
	return dto
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
