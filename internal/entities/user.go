package entities

// UserID defines the unique user's identifier.
type UserID string

// User defines the member of the Team.
type User struct {
	ID       UserID `json:"user_id"`
	Name     string `json:"username"`
	IsActive bool   `json:"is_active"`
	TeamName string `json:"team_name"`
}
