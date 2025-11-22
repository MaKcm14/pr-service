package dto

import "github.com/MaKcm14/pr-service/internal/entities"

// TeamMember defines the dto object for the User's team view.
type TeamMember struct {
	ID       entities.UserID `json:"user_id"`
	Name     string          `json:"username"`
	IsActive bool            `json:"is_active"`
}

// TeamDTO defines the dto object for the Team's view.
type TeamDTO struct {
	Name    string       `json:"team_name"`
	Members []TeamMember `json:"members"`
}

func NewTeamDTO() TeamDTO {
	return TeamDTO{
		Members: make([]TeamMember, 0, 250),
	}
}

func UserToTeamMember(user entities.User) TeamMember {
	return TeamMember{
		ID:       user.ID,
		Name:     user.Name,
		IsActive: user.IsActive,
	}
}

func TeamToTeamDTO(team entities.Team) TeamDTO {
	dto := NewTeamDTO()

	dto.Name = team.Name
	for _, member := range team.Members {
		dto.Members = append(dto.Members, UserToTeamMember(member))
	}
	return dto
}
