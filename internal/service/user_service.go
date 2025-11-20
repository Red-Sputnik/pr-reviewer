package service

import (
	"errors"
	"pr-reviewer/internal/model"
	"pr-reviewer/internal/repository"
)

var ErrUserNotFound = errors.New("user not found")

type UserService struct {
	userRepo *repository.UserRepo
	teamRepo *repository.TeamRepo
}

func NewUserService(userRepo *repository.UserRepo, teamRepo *repository.TeamRepo) *UserService {
	return &UserService{
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}

// SetIsActive устанавливает флаг активности пользователя
func (s *UserService) SetIsActive(userID string, isActive bool) (*model.User, error) {
	user, err := s.userRepo.SetIsActive(userID, isActive)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

// GetUserByID возвращает пользователя по ID
func (s *UserService) GetUserByID(userID string) (*model.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

// GetUsersByTeam возвращает всех пользователей команды
func (s *UserService) GetUsersByTeam(teamName string) ([]model.User, error) {
	return s.userRepo.GetByTeam(teamName)
}
