package service

import (
	"errors"
	"fmt"
	"strconv"

	"pr-reviewer/internal/model"
	"pr-reviewer/internal/repository"
)

// Ошибки доменной логики
var (
	ErrTeamExists   = errors.New("team already exists")
	ErrTeamNotFound = errors.New("team not found")
)

// TeamService управляет командами
type TeamService struct {
	teamRepo *repository.TeamRepo
	userRepo *repository.UserRepo
}

// NewTeamService создаёт новый TeamService
func NewTeamService(teamRepo *repository.TeamRepo, userRepo *repository.UserRepo) *TeamService {
	return &TeamService{
		teamRepo: teamRepo,
		userRepo: userRepo,
	}
}

// parseID конвертирует строковый UserID в int64 с проверкой ошибки
func parseID(id string) (int64, error) {
	return strconv.ParseInt(id, 10, 64)
}

// CreateOrUpdateTeam создаёт команду или обновляет участников
func (s *TeamService) CreateOrUpdateTeam(team *model.Team) (*model.Team, error) {
	existingTeam, err := s.teamRepo.GetByName(team.Name)
	if err != nil && err != repository.ErrNotFound {
		return nil, err
	}

	if existingTeam != nil {
		// команда уже существует, обновляем участников
		for i := range team.Members {
			id, err := parseID(team.Members[i].UserID)
			if err != nil {
				return nil, fmt.Errorf("invalid user ID %s: %w", team.Members[i].UserID, err)
			}

			u := &model.User{
				ID:       id,
				Username: team.Members[i].Username,
				IsActive: team.Members[i].IsActive,
			}

			if err := s.userRepo.CreateOrUpdate(u); err != nil {
				return nil, err
			}
		}
		return existingTeam, ErrTeamExists
	}

	// создаём команду
	if err := s.teamRepo.Create(team); err != nil {
		return nil, err
	}

	// создаём пользователей в команде
	for i := range team.Members {
		id, err := parseID(team.Members[i].UserID)
		if err != nil {
			return nil, fmt.Errorf("invalid user ID %s: %w", team.Members[i].UserID, err)
		}

		u := &model.User{
			ID:       id,
			Username: team.Members[i].Username,
			IsActive: team.Members[i].IsActive,
		}

		if err := s.userRepo.CreateOrUpdate(u); err != nil {
			return nil, err
		}
	}

	return team, nil
}

// GetTeam возвращает команду по имени
func (s *TeamService) GetTeam(teamName string) (*model.Team, error) {
	team, err := s.teamRepo.GetByName(teamName)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrTeamNotFound
		}
		return nil, err
	}

	// загружаем участников
	users, err := s.userRepo.GetByTeam(teamName) // []User
	if err != nil {
		return nil, err
	}

	// конвертируем []User -> []TeamMember
	members := make([]model.TeamMember, len(users))
	for i, u := range users {
		members[i] = model.TeamMember{
			UserID:   fmt.Sprintf("%d", u.ID),
			Username: u.Username,
			IsActive: u.IsActive,
		}
	}

	team.Members = members
	return team, nil
}
