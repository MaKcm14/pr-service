package entities

// TeamID defines the unique user's identifier.
type TeamID int64

// Team defines the groups of the 'User's.
type Team struct {
	ID      TeamID `json:"-"`
	Name    string `json:"team_name"`
	Members []User `json:"members"`
}

func NewTeam() Team {
	return Team{
		Members: make([]User, 250),
	}
}
