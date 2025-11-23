package dto

import (
	"time"

	"github.com/MaKcm14/pr-service/internal/entities"
)

type PullRequestDTO struct {
	ID        entities.PullRequestID     `json:"pull_request_id"`
	Name      string                     `json:"pull_request_name"`
	Status    entities.PullRequestStatus `json:"status"`
	CreatedAt *time.Time                 `json:"created_at"`
	MergedAt  *time.Time                 `json:"merged_at"`
	AuthorID  entities.UserID            `json:"author_id"`
	Reviewers []entities.UserID          `json:"assigned_reviewers"`
}

func NewPullRequestDTO() PullRequestDTO {
	return PullRequestDTO{
		Reviewers: make([]entities.UserID, 0, 5),
	}
}

func PullRequestToPullRequestDTO(pullReq entities.PullRequest) PullRequestDTO {
	dto := PullRequestDTO{
		ID:        pullReq.ID,
		Name:      pullReq.Name,
		Status:    pullReq.Status,
		AuthorID:  pullReq.Author.ID,
		Reviewers: make([]entities.UserID, 0, len(pullReq.Reviewers)),
	}

	if pullReq.CreatedAt != nil {
		dto.CreatedAt = new(time.Time)
		(*dto.CreatedAt) = (*pullReq.CreatedAt)
	}

	if pullReq.MergedAt != nil {
		dto.MergedAt = new(time.Time)
		(*dto.MergedAt) = (*pullReq.MergedAt)
	}

	for _, user := range pullReq.Reviewers {
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

func PullRequestDTOToPullRequest(pullReq PullRequestDTO) entities.PullRequest {
	obj := entities.PullRequest{
		ID:     pullReq.ID,
		Name:   pullReq.Name,
		Status: pullReq.Status,
		Author: entities.User{
			ID: pullReq.AuthorID,
		},
		Reviewers: make(map[entities.UserID]entities.User, len(pullReq.Reviewers)),
	}
	if pullReq.CreatedAt != nil {
		obj.CreatedAt = new(time.Time)
		(*obj.CreatedAt) = (*pullReq.CreatedAt)
	}

	if pullReq.MergedAt != nil {
		obj.MergedAt = new(time.Time)
		(*obj.MergedAt) = (*pullReq.MergedAt)
	}

	for _, userID := range pullReq.Reviewers {
		obj.Reviewers[userID] = entities.User{
			ID: userID,
		}
	}

	return obj
}

type PullRequestChangeReviewerDTO struct {
	ID            entities.PullRequestID `json:"pull_request_id"`
	OldReviewerID entities.UserID        `json:"old_reviewer_id"`
}
