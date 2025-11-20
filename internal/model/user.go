package model

// User описывает пользователя
type User struct {
	ID       int64  `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}
