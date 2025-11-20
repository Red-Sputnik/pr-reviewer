package service

import (
	"errors"
	"math/rand"
	"pr-reviewer/internal/model"
	"pr-reviewer/internal/repository"
	"strconv"
	"time"
)

var (
	ErrPRNotFound          = errors.New("PR not found")
	ErrPRAlreadyMerged     = errors.New("PR already merged")
	ErrNoCandidate         = errors.New("no active replacement candidate in team")
	ErrReviewerNotAssigned = errors.New("reviewer is not assigned to this PR")
)

type PRService struct {
	prRepo   *repository.PRRepo
	userRepo *repository.UserRepo
	teamRepo *repository.TeamRepo
}

func NewPRService(prRepo *repository.PRRepo, userRepo *repository.UserRepo, teamRepo *repository.TeamRepo) *PRService {
	return &PRService{
		prRepo:   prRepo,
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}

// CreatePR создает PR и назначает до 2 ревьюверов
func (s *PRService) CreatePR(pr *model.PullRequest) (*model.PullRequest, error) {
	existing, _ := s.prRepo.GetByID(pr.ID)
	if existing != nil {
		return nil, errors.New("PR_EXISTS")
	}

	pr.Status = "OPEN"
	pr.CreatedAt = time.Now()

	// получаем команду автора
	user, err := s.userRepo.GetByID(pr.AuthorID)
	if err != nil {
		return nil, ErrPRNotFound
	}
	users, _ := s.userRepo.GetByTeam(user.TeamName)

	var reviewers []string
	for _, u := range users {
		userIDStr := strconv.FormatInt(u.ID, 10)
		if userIDStr != pr.AuthorID && u.IsActive {
			reviewers = append(reviewers, userIDStr)
		}
	}
	// выбираем случайно до 2 ревьюверов
	rand.Shuffle(len(reviewers), func(i, j int) { reviewers[i], reviewers[j] = reviewers[j], reviewers[i] })
	if len(reviewers) > 2 {
		reviewers = reviewers[:2]
	}
	pr.AssignedReviewers = reviewers

	if err := s.prRepo.Create(pr); err != nil {
		return nil, err
	}
	return pr, nil
}

// MergePR помечает PR как MERGED (идемпотентно)
func (s *PRService) MergePR(prID string) (*model.PullRequest, error) {
	pr, err := s.prRepo.GetByID(prID)
	if err != nil {
		return nil, ErrPRNotFound
	}
	if pr.Status == "MERGED" {
		return pr, nil
	}
	now := time.Now()
	pr.Status = "MERGED"
	pr.MergedAt = &now
	if err := s.prRepo.Update(pr); err != nil {
		return nil, err
	}
	return pr, nil
}

// ReassignReviewer заменяет ревьювера на случайного активного пользователя из команды
func (s *PRService) ReassignReviewer(prID, oldUserID string) (*model.PullRequest, string, error) {
	pr, err := s.prRepo.GetByID(prID)
	if err != nil {
		return nil, "", ErrPRNotFound
	}
	if pr.Status == "MERGED" {
		return nil, "", ErrPRAlreadyMerged
	}

	// проверяем что oldUserID назначен
	found := false
	for _, r := range pr.AssignedReviewers {
		if r == oldUserID {
			found = true
			break
		}
	}
	if !found {
		return nil, "", ErrReviewerNotAssigned
	}

	oldUser, err := s.userRepo.GetByID(oldUserID)
	if err != nil {
		return nil, "", err
	}
	teamUsers, _ := s.userRepo.GetByTeam(oldUser.TeamName)

	var candidates []string
	for _, u := range teamUsers {
		userIDStr := strconv.FormatInt(u.ID, 10)
		if userIDStr != oldUserID && u.IsActive {
			candidates = append(candidates, userIDStr)
		}
	}
	if len(candidates) == 0 {
		return nil, "", ErrNoCandidate
	}

	rand.Shuffle(len(candidates), func(i, j int) { candidates[i], candidates[j] = candidates[j], candidates[i] })
	newReviewer := candidates[0]

	for i, r := range pr.AssignedReviewers {
		if r == oldUserID {
			pr.AssignedReviewers[i] = newReviewer
			break
		}
	}

	if err := s.prRepo.Update(pr); err != nil {
		return nil, "", err
	}
	return pr, newReviewer, nil
}

// GetPRsForReviewer возвращает PR, где пользователь назначен ревьювером
func (s *PRService) GetPRsForReviewer(userID string) ([]model.PullRequest, error) {
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return s.prRepo.GetByReviewer(userID)
}
