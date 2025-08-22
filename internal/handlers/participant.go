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

func (h *ParticipantHandler) HandleParticipants(w http.ResponseWriter, r *http.Request) {
	log.Printf("Method: %s, URL: %s", r.Method, r.URL.String())  // デバッグ用ログ

	switch r.Method {
	case "GET":
		roomID := r.URL.Query().Get("room_id")
		if roomID == "" {
			http.Error(w, "room_idは必須です", http.StatusBadRequest)
			return
		}
		h.GetParticipantsByRoomID(w, r, roomID)
	case "POST":
		h.AddParticipant(w, r)
	default:
		http.Error(w, "サポートされていないメソッドです", http.StatusMethodNotAllowed)
	}
}



func (h *ParticipantHandler) GetParticipantsByRoomID(w http.ResponseWriter, r *http.Request, roomID string) {
	// 部屋の存在確認
	var roomExists bool
	err := h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM rooms WHERE id = $1)", roomID).Scan(&roomExists)
	if err != nil {
		log.Println("部屋の存在確認に失敗しました:", err)
		http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
		return
	}
	if !roomExists {
		http.Error(w, "指定された部屋は存在しません", http.StatusNotFound)
		return
	}

	// 参加者一覧を取得
	rows, err := h.db.Query(`
		SELECT u.id, u.user_name
		FROM participants p
		JOIN users u ON p.user_id = u.id
		WHERE p.room_id = $1
		ORDER BY u.user_name
	`, roomID)
	if err != nil {
		log.Println("参加者の取得に失敗しました:", err)
		http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.ParticipantUser
	for rows.Next() {
		var user models.ParticipantUser
		err := rows.Scan(&user.ID, &user.Name)
		if err != nil {
			log.Println("参加者データの読み取りに失敗しました:", err)
			continue
		}
		users = append(users, user)
	}

	// レスポンスを作成
	response := models.ParticipantsResponse{
		RoomID: roomID,
		Users:  users,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("レスポンスの作成に失敗しました:", err)
	}
}

func (h *ParticipantHandler) AddParticipant(w http.ResponseWriter, r *http.Request) {
	var req models.ParticipantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "リクエストボディが不正です", http.StatusBadRequest)
		return
	}

	// バリデーション
	if req.RoomID == "" || req.UserID == "" {
		http.Error(w, "room_idとuser_idは必須です", http.StatusBadRequest)
		return
	}

	// 部屋の存在確認
	var roomExists bool
	err := h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM rooms WHERE id = $1)", req.RoomID).Scan(&roomExists)
	if err != nil {
		log.Println("部屋の存在確認に失敗しました:", err)
		http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
		return
	}
	if !roomExists {
		http.Error(w, "指定された部屋は存在しません", http.StatusNotFound)
		return
	}

	// ユーザーの存在確認
	var userExists bool
	err = h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", req.UserID).Scan(&userExists)
	if err != nil {
		log.Println("ユーザーの存在確認に失敗しました:", err)
		http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
		return
	}
	if !userExists {
		http.Error(w, "指定されたユーザーは存在しません", http.StatusNotFound)
		return
	}

	// 参加者の追加
	_, err = h.db.Exec(
		"INSERT INTO participants (room_id, user_id) VALUES ($1, $2) ON CONFLICT (room_id, user_id) DO NOTHING",
		req.RoomID,
		req.UserID,
	)
	if err != nil {
		log.Println("参加者の追加に失敗しました:", err)
		http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
		return
	}

	// レスポンスの作成
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{
		"room_id": req.RoomID,
		"user_id": req.UserID,
	}); err != nil {
		log.Println("レスポンスの作成に失敗しました:", err)
	}
}
