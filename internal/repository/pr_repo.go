package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"pr-reviewer/internal/model"
)

type PRRepo struct {
	db *sql.DB
}

func NewPRRepo(db *sql.DB) *PRRepo {
	return &PRRepo{db: db}
}

// Create создает PR
func (r *PRRepo) Create(pr *model.PullRequest) error {
	reviewersJSON, _ := json.Marshal(pr.AssignedReviewers)
	query := `
	INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status, assigned_reviewers, created_at)
	VALUES ($1,$2,$3,$4,$5,$6)
	`
	_, err := r.db.Exec(query, pr.ID, pr.Name, pr.AuthorID, pr.Status, reviewersJSON, pr.CreatedAt)
	return err
}

// GetByID возвращает PR по ID
func (r *PRRepo) GetByID(prID string) (*model.PullRequest, error) {
	query := `SELECT pull_request_id, pull_request_name, author_id, status, assigned_reviewers, created_at, merged_at FROM pull_requests WHERE pull_request_id=$1`
	row := r.db.QueryRow(query, prID)

	var pr model.PullRequest
	var reviewersJSON []byte
	if err := row.Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status, &reviewersJSON, &pr.CreatedAt, &pr.MergedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	json.Unmarshal(reviewersJSON, &pr.AssignedReviewers)
	return &pr, nil
}

// Update обновляет PR
func (r *PRRepo) Update(pr *model.PullRequest) error {
	reviewersJSON, _ := json.Marshal(pr.AssignedReviewers)
	query := `
	UPDATE pull_requests
	SET status=$1, assigned_reviewers=$2, merged_at=$3
	WHERE pull_request_id=$4
	`
	_, err := r.db.Exec(query, pr.Status, reviewersJSON, pr.MergedAt, pr.ID)
	return err
}

// GetByReviewer возвращает PR, где пользователь ревьювер
func (r *PRRepo) GetByReviewer(userID string) ([]model.PullRequest, error) {
	query := `SELECT pull_request_id, pull_request_name, author_id, status, assigned_reviewers, created_at, merged_at FROM pull_requests`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []model.PullRequest
	for rows.Next() {
		var pr model.PullRequest
		var reviewersJSON []byte
		if err := rows.Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status, &reviewersJSON, &pr.CreatedAt, &pr.MergedAt); err != nil {
			return nil, err
		}
		json.Unmarshal(reviewersJSON, &pr.AssignedReviewers)
		for _, rID := range pr.AssignedReviewers {
			if rID == userID {
				prs = append(prs, pr)
				break
			}
		}
	}
	return prs, nil
}
