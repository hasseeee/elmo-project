package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/shuto.sawaki/elmo-project/internal/ai"
	"github.com/shuto.sawaki/elmo-project/internal/models"
)

// ... (GetRooms, CreateRoom, GetRoomByID, SaveConclusion は変更なし) ...

type RoomHandler struct {
	db          *sql.DB
	aiGenerator ai.AIGenerator
}

func NewRoomHandler(db *sql.DB, aiGen ai.AIGenerator) *RoomHandler {
	return &RoomHandler{
		db:          db,
		aiGenerator: aiGen,
	}
}

// GET /rooms
func (h *RoomHandler) GetRooms(w http.ResponseWriter, r *http.Request) {
	// ... (変更なし)
}

// POST /rooms
func (h *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	// ... (変更なし)
}

// GET /rooms/{id}
func (h *RoomHandler) GetRoomByID(w http.ResponseWriter, r *http.Request) {
	// ... (変更なし)
}

// POST /rooms/{id}/conclusion
func (h *RoomHandler) SaveConclusion(w http.ResponseWriter, r *http.Request) {
	// ... (変更なし)
}

// GET /rooms/{id}/start
func (h *RoomHandler) StartRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "GETメソッドのみサポートされています", http.StatusMethodNotAllowed)
		return
	}
	log.Println("StartRoom: リクエストを受信しました")
	roomID := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/rooms/"), "/start")

	var room models.Room
	err := h.db.QueryRow("SELECT id, title, description, status FROM rooms WHERE id = $1", roomID).Scan(&room.ID, &room.Title, &room.Description, &room.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "指定された部屋は見つかりません", http.StatusNotFound)
		} else {
			log.Println("データベースクエリの実行に失敗しました:", err)
			http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
		}
		return
	}

	if room.Status != "not started" {
		http.Error(w, "部屋は既に開始されています", http.StatusBadRequest)
		return
	}

	initialQuestion, err := h.aiGenerator.GenerateInitialQuestion(r.Context(), room.Title, room.Description)
	if err != nil {
		log.Println("AIからの問いかけ生成に失敗しました:", err)
		http.Error(w, "AI API呼び出しエラー", http.StatusInternalServerError)
		return
	}

	_, err = h.db.Exec("UPDATE rooms SET status = $1, initial_question = $2 WHERE id = $3", "inprogress", initialQuestion, roomID)
	if err != nil {
		log.Println("部屋の状態更新に失敗しました:", err)
		http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
		return
	}

	// ★★★ ここから修正 ★★★
	rows, err := h.db.Query(`SELECT u.id, u.user_name FROM participants p JOIN users u ON p.user_id = u.id WHERE p.room_id = $1`, roomID)
	if err != nil {
		log.Println("参加者リストの取得に失敗しました:", err)
		http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var participants []models.ParticipantUser
	for rows.Next() {
		var p models.ParticipantUser
		if err := rows.Scan(&p.ID, &p.Name); err != nil {
			log.Println("参加者情報の読み取りに失敗しました:", err)
			http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
			return
		}
		participants = append(participants, p)
	}

	response := models.StartRoomResponse{
		InitialQuestion: initialQuestion,
		RoomInfo: models.RoomInfo{
			RoomID: roomID,
			Title:  room.Title,
			Status: "inprogress",
		},
		Participants: participants,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("JSONへの変換に失敗しました:", err)
	}
	// ★★★ ここまで修正 ★★★
}

// POST /rooms/{id}/sorena
func (h *RoomHandler) HandleSorena(w http.ResponseWriter, r *http.Request) {
    // ... (sorena.goから持ってきたロジック。変更なし)
}