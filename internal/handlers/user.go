package handlers

import (
	"encoding/json"
	"net/http"
	"pr-reviewer/internal/service"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	userService *service.UserService
	prService   *service.PRService
}

func NewUserHandler(userService *service.UserService, prService *service.PRService) *UserHandler {
	return &UserHandler{
		userService: userService,
		prService:   prService,
	}
}

func (h *UserHandler) RegisterUserRoutes(r *mux.Router) {
	r.HandleFunc("/users/setIsActive", h.SetIsActive).Methods("POST")
	r.HandleFunc("/users/getReview", h.GetReviewPRs).Methods("GET")
}

// SetIsActive обновляет флаг активности
func (h *UserHandler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID   string `json:"user_id"`
		IsActive bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": {"code":"INVALID_REQUEST","message":"invalid JSON"}}`, http.StatusBadRequest)
		return
	}

	user, err := h.userService.SetIsActive(req.UserID, req.IsActive)
	if err != nil {
		if err == service.ErrUserNotFound {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]string{
					"code":    "NOT_FOUND",
					"message": "user not found",
				},
			})
			return
		}
		http.Error(w, `{"error":{"code":"INTERNAL","message":"internal error"}}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"user": user,
	})
}

// GetReviewPRs возвращает PR, где пользователь назначен ревьювером
func (h *UserHandler) GetReviewPRs(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, `{"error":{"code":"INVALID_REQUEST","message":"user_id required"}}`, http.StatusBadRequest)
		return
	}

	prs, err := h.prService.GetPRsForReviewer(userID)
	if err != nil {
		if err == service.ErrUserNotFound {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]string{
					"code":    "NOT_FOUND",
					"message": "user not found",
				},
			})
			return
		}
		http.Error(w, `{"error":{"code":"INTERNAL","message":"internal error"}}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":       userID,
		"pull_requests": prs,
	})
}
