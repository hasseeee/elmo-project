package handlers

import (
	"database/sql"
	"encoding/json"
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

	rows, err := h.db.Query("SELECT u.id, u.user_name FROM participants p JOIN users u ON p.user_id = u.id WHERE p.room_id = $1", roomID)
	if err != nil {
		http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.ParticipantUser
	for rows.Next() {
		var user models.ParticipantUser
		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	response := models.ParticipantsResponse{
		RoomID: roomID,
		Users:  users,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// POST /participants
func (h *ParticipantHandler) AddParticipant(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POSTメソッドのみサポートされています", http.StatusMethodNotAllowed)
		return
	}

	var req models.ParticipantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "リクエストボディが不正です", http.StatusBadRequest)
		return
	}

	if req.RoomID == "" || req.UserID == "" {
		http.Error(w, "room_idとuser_idは必須です", http.StatusBadRequest)
		return
	}

	_, err := h.db.Exec("INSERT INTO participants (room_id, user_id) VALUES ($1, $2)", req.RoomID, req.UserID)
	if err != nil {
		http.Error(w, "参加者の追加に失敗しました", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}