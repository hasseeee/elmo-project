package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/shuto.sawaki/elmo-project/internal/models"
)

type ParticipantHandler struct {
	db *sql.DB
}

func NewParticipantHandler(db *sql.DB) *ParticipantHandler {
	return &ParticipantHandler{db: db}
}

// GET /participants
func (h *ParticipantHandler) GetParticipants(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "GETメソッドのみサポートされています", http.StatusMethodNotAllowed)
		return
	}

	roomID := r.URL.Query().Get("room_id")
	if roomID == "" {
		http.Error(w, "room_idは必須です", http.StatusBadRequest)
		return
	}
	
	// ... (以降のGetParticipantsByRoomIDのロジック)
}

// POST /participants
func (h *ParticipantHandler) AddParticipant(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POSTメソッドのみサポートされています", http.StatusMethodNotAllowed)
		return
	}

	var req models.ParticipantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// ... (以降のAddParticipantのロジック)
	}
}