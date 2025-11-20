package handlers

import (
	"encoding/json"
	"net/http"
	"pr-reviewer/internal/model"
	"pr-reviewer/internal/service"

	"github.com/gorilla/mux"
)

type PRHandler struct {
	prService *service.PRService
}

func NewPRHandler(prService *service.PRService) *PRHandler {
	return &PRHandler{prService: prService}
}

func (h *PRHandler) RegisterPRRoutes(r *mux.Router) {
	r.HandleFunc("/pullRequest/create", h.CreatePR).Methods("POST")
	r.HandleFunc("/pullRequest/merge", h.MergePR).Methods("POST")
	r.HandleFunc("/pullRequest/reassign", h.ReassignReviewer).Methods("POST")
}

// CreatePR создаёт PR и назначает ревьюверов
func (h *PRHandler) CreatePR(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID     string `json:"pull_request_id"`
		Name   string `json:"pull_request_name"`
		Author string `json:"author_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":{"code":"INVALID_REQUEST","message":"invalid JSON"}}`, http.StatusBadRequest)
		return
	}

	pr, err := h.prService.CreatePR(&model.PullRequest{
		ID:       req.ID,
		Name:     req.Name,
		AuthorID: req.Author,
	})
	if err != nil {
		if err.Error() == "PR_EXISTS" {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]string{"code": "PR_EXISTS", "message": "PR id already exists"},
			})
			return
		}
		http.Error(w, `{"error":{"code":"INTERNAL","message":"internal error"}}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"pr": pr})
}

// MergePR помечает PR как MERGED
func (h *PRHandler) MergePR(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID string `json:"pull_request_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":{"code":"INVALID_REQUEST","message":"invalid JSON"}}`, http.StatusBadRequest)
		return
	}

	pr, err := h.prService.MergePR(req.ID)
	if err != nil {
		if err == service.ErrPRNotFound {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]string{"code": "NOT_FOUND", "message": "PR not found"},
			})
			return
		}
		http.Error(w, `{"error":{"code":"INTERNAL","message":"internal error"}}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"pr": pr})
}

// ReassignReviewer переназначает ревьювера
func (h *PRHandler) ReassignReviewer(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PRID      string `json:"pull_request_id"`
		OldUserID string `json:"old_user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":{"code":"INVALID_REQUEST","message":"invalid JSON"}}`, http.StatusBadRequest)
		return
	}

	pr, newID, err := h.prService.ReassignReviewer(req.PRID, req.OldUserID)
	if err != nil {
		switch err {
		case service.ErrPRNotFound:
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]interface{}{"error": map[string]string{"code": "NOT_FOUND", "message": "PR or user not found"}})
		case service.ErrPRAlreadyMerged:
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]interface{}{"error": map[string]string{"code": "PR_MERGED", "message": "cannot reassign on merged PR"}})
		case service.ErrReviewerNotAssigned:
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]interface{}{"error": map[string]string{"code": "NOT_ASSIGNED", "message": "reviewer is not assigned to this PR"}})
		case service.ErrNoCandidate:
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]interface{}{"error": map[string]string{"code": "NO_CANDIDATE", "message": "no active replacement candidate in team"}})
		default:
			http.Error(w, `{"error":{"code":"INTERNAL","message":"internal error"}}`, http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"pr": pr, "replaced_by": newID})
}
