package handlers

import (
	"encoding/json"
	"net/http"

	"pr-reviewer/internal/model"
	"pr-reviewer/internal/service"

	"github.com/gorilla/mux"
)

// TeamHandler объединяет сервис команд
type TeamHandler struct {
	teamService *service.TeamService
}

// NewTeamHandler создаёт новый обработчик команд
func NewTeamHandler(teamService *service.TeamService) *TeamHandler {
	return &TeamHandler{teamService: teamService}
}

// RegisterTeamRoutes регистрирует маршруты команд
func (h *TeamHandler) RegisterTeamRoutes(r *mux.Router) {
	r.HandleFunc("/team/add", h.AddTeam).Methods("POST")
	r.HandleFunc("/team/get", h.GetTeam).Methods("GET")
}

// AddTeam создаёт команду с участниками
func (h *TeamHandler) AddTeam(w http.ResponseWriter, r *http.Request) {
	var team model.Team
	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		http.Error(w, `{"error": {"code": "INVALID_REQUEST", "message": "invalid JSON"}}`, http.StatusBadRequest)
		return
	}

	createdTeam, err := h.teamService.CreateOrUpdateTeam(&team)
	if err != nil {
		if err == service.ErrTeamExists {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]string{
					"code":    "TEAM_EXISTS",
					"message": "team_name already exists",
				},
			})
			return
		}
		http.Error(w, `{"error": {"code": "INTERNAL", "message": "internal error"}}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"team": createdTeam,
	})
}

// GetTeam возвращает команду по имени
func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		http.Error(w, `{"error": {"code": "INVALID_REQUEST", "message": "team_name query required"}}`, http.StatusBadRequest)
		return
	}

	team, err := h.teamService.GetTeam(teamName)
	if err != nil {
		if err == service.ErrTeamNotFound {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]string{
					"code":    "NOT_FOUND",
					"message": "team not found",
				},
			})
			return
		}
		http.Error(w, `{"error": {"code": "INTERNAL", "message": "internal error"}}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(team)
}
