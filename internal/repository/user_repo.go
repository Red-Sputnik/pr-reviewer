package repository

import (
	"database/sql"
	"errors"
	"pr-reviewer/internal/model"
)

// UserRepo работает с пользователями
type UserRepo struct {
	db *sql.DB
}

// NewUserRepo создаёт новый UserRepo
func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

// CreateOrUpdate создаёт или обновляет пользователя
func (r *UserRepo) CreateOrUpdate(user *model.User) error {
	query := `
	INSERT INTO users (user_id, username, team_name, is_active)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (user_id) DO UPDATE SET username=$2, team_name=$3, is_active=$4
	`
	_, err := r.db.Exec(query, user.ID, user.Username, user.TeamName, user.IsActive)
	return err
}

// GetByID возвращает пользователя по ID
func (r *UserRepo) GetByID(userID string) (*model.User, error) {
	query := `SELECT user_id, username, team_name, is_active FROM users WHERE user_id=$1`
	row := r.db.QueryRow(query, userID)

	var u model.User
	if err := row.Scan(&u.ID, &u.Username, &u.TeamName, &u.IsActive); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

// GetByTeam возвращает всех пользователей команды
func (r *UserRepo) GetByTeam(teamName string) ([]model.User, error) {
	query := `SELECT user_id, username, team_name, is_active FROM users WHERE team_name=$1`
	rows, err := r.db.Query(query, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Username, &u.TeamName, &u.IsActive); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// SetIsActive обновляет флаг активности пользователя
func (r *UserRepo) SetIsActive(userID string, isActive bool) (*model.User, error) {
	query := `UPDATE users SET is_active=$1 WHERE user_id=$2 RETURNING user_id, username, team_name, is_active`
	row := r.db.QueryRow(query, isActive, userID)

	var u model.User
	if err := row.Scan(&u.ID, &u.Username, &u.TeamName, &u.IsActive); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}
