package model

// TeamMember описывает участника команды
type TeamMember struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

// Team описывает команду
type Team struct {
	Name    string       `json:"team_name"`
	Members []TeamMember `json:"members"`
}
