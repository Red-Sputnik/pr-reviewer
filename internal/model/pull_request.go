package model

import "time"

type PullRequest struct {
	ID                string     `db:"pull_request_id" json:"pull_request_id"`
	Name              string     `db:"pull_request_name" json:"pull_request_name"`
	AuthorID          string     `db:"author_id" json:"author_id"`
	Status            string     `db:"status" json:"status"` // OPEN|MERGED
	AssignedReviewers []string   `db:"assigned_reviewers" json:"assigned_reviewers"`
	CreatedAt         time.Time  `db:"created_at" json:"createdAt"`
	MergedAt          *time.Time `db:"merged_at" json:"mergedAt,omitempty"`
}
