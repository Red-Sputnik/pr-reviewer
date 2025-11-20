package repository

import (
	"database/sql"
	"errors"
	"pr-reviewer/internal/model"
)

// Ошибки репозитория
var ErrNotFound = errors.New("record not found")

// TeamRepo работает с командами
type TeamRepo struct {
	db *sql.DB
}

// NewTeamRepo создаёт новый TeamRepo
func NewTeamRepo(db *sql.DB) *TeamRepo {
	return &TeamRepo{db: db}
}

// Create создаёт новую команду
func (r *TeamRepo) Create(team *model.Team) error {
	query := `INSERT INTO teams (name) VALUES ($1)`
	_, err := r.db.Exec(query, team.Name)
	return err
}

// GetByName возвращает команду по имени
func (r *TeamRepo) GetByName(name string) (*model.Team, error) {
	query := `SELECT name FROM teams WHERE name=$1`
	row := r.db.QueryRow(query, name)

	var team model.Team
	if err := row.Scan(&team.Name); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &team, nil
}

// Delete удаляет команду (для тестов или администрирования)
func (r *TeamRepo) Delete(name string) error {
	query := `DELETE FROM teams WHERE name=$1`
	_, err := r.db.Exec(query, name)
	return err
}
