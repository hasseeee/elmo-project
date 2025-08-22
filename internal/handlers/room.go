package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/matoous/go-nanoid/v2"
	"github.com/shuto.sawaki/elmo-project/internal/ai"
	"github.com/shuto.sawaki/elmo-project/internal/models"
)

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
	if r.Method != http.MethodGet {
		http.Error(w, "GETメソッドのみサポートされています", http.StatusMethodNotAllowed)
		return
	}

	log.Println("GetRooms: リクエストを受信しました")
	rows, err := h.db.Query("SELECT id, title, description FROM rooms ORDER BY id ASC")
	if err != nil {
		log.Println("データベースクエリの実行に失敗しました:", err)
		http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var rooms []models.Room
	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.Title, &room.Description)
		if err != nil {
			log.Println("データベースからのデータ読み取りに失敗しました:", err)
			http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
			return
		}
		rooms = append(rooms, room)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rooms); err != nil {
		log.Println("JSONへの変換に失敗しました:", err)
	}
}

// POST /rooms
func (h *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POSTメソッドのみサポートされています", http.StatusMethodNotAllowed)
		return
	}

	var newRoom models.Room
	if err := json.NewDecoder(r.Body).Decode(&newRoom); err != nil {
		http.Error(w, "リクエストボディが不正です", http.StatusBadRequest)
		return
	}

	if newRoom.Title == "" {
		http.Error(w, "タイトルは必須です", http.StatusBadRequest)
		return
	}

	newId, err := gonanoid.Generate("0123456789abcdefghijklmnopqrstuvwxyz", 6)
	if err != nil {
		log.Println("IDの生成に失敗しました:", err)
		http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
		return
	}
	newRoom.ID = newId

	sqlStatement := `INSERT INTO rooms (id, title, description) VALUES ($1, $2, $3)`
	_, err = h.db.Exec(sqlStatement, newRoom.ID, newRoom.Title, newRoom.Description)
	if err != nil {
		log.Println("データベースへのINSERTに失敗しました:", err)
		http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
		return
	}

	log.Printf("新しい部屋を作成しました: ID=%s, Title=%s", newRoom.ID, newRoom.Title)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newRoom)
}

// GET /rooms/{id}
func (h *RoomHandler) GetRoomByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "GETメソッドのみサポートされています", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/rooms/")
	var room models.Room
	sqlStatement := `SELECT id, title, description, conclusion, status, initial_question FROM rooms WHERE id = $1`
	err := h.db.QueryRow(sqlStatement, id).Scan(&room.ID, &room.Title, &room.Description, &room.Conclusion, &room.Status, &room.InitialQuestion)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "指定された部屋は見つかりません", http.StatusNotFound)
		} else {
			log.Println("データベースクエリの実行に失敗しました:", err)
			http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(room)
}

// POST /rooms/{id}/conclusion
func (h *RoomHandler) SaveConclusion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POSTメソッドのみサポートされています", http.StatusMethodNotAllowed)
		return
	}

	roomID := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/rooms/"), "/conclusion")

	var req models.ConclusionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "リクエストボディが不正です", http.StatusBadRequest)
		return
	}

	if req.Conclusion == "" {
		http.Error(w, "結論は必須です", http.StatusBadRequest)
		return
	}

	// 結論を保存
	_, err := h.db.Exec("UPDATE rooms SET conclusion = $1, status = 'concluded' WHERE id = $2", req.Conclusion, roomID)
	if err != nil {
		log.Println("結論の保存に失敗しました:", err)
		http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
		return
	}

	// 更新後の部屋の情報を取得して返す
	var room models.Room
	err = h.db.QueryRow("SELECT id, title, description, conclusion, status FROM rooms WHERE id = $1", roomID).
		Scan(&room.ID, &room.Title, &room.Description, &room.Conclusion, &room.Status)
	if err != nil {
		log.Println("更新後の部屋情報の取得に失敗しました:", err)
		http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(room)
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
		// ... (エラー処理は省略)
	}

	if room.Status != "not started" {
		http.Error(w, "部屋は既に開始されています", http.StatusBadRequest)
		return
	}

	// AIを使って最初の問いかけを生成
	initialQuestion, err := h.aiGenerator.GenerateInitialQuestion(r.Context(), room.Title, room.Description)
	if err != nil {
		log.Println("AIからの問いかけ生成に失敗しました:", err)
		http.Error(w, "AI API呼び出しエラー", http.StatusInternalServerError)
		return
	}

	// 部屋の状態を更新
	_, err = h.db.Exec("UPDATE rooms SET status = $1, initial_question = $2 WHERE id = $3", "inprogress", initialQuestion, roomID)
	if err != nil {
		// ... (エラー処理は省略)
	}

	// 参加者リストを取得
	rows, err := h.db.Query(`SELECT u.id, u.user_name FROM participants p JOIN users u ON p.user_id = u.id WHERE p.room_id = $1`, roomID)
	// ... (参加者リストの取得とレスポンス生成処理は省略)

	// レスポンスを返す
	// ... (省略)
}

// POST /rooms/{id}/sorena
func (h *RoomHandler) HandleSorena(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POSTメソッドのみサポートされています", http.StatusMethodNotAllowed)
		return
	}
	roomID := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/rooms/"), "/sorena")

	var req models.SorenaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// ... (エラー処理は省略)
	}
	// ... (以降の処理はsorena.goから持ってくる)
	sqlStatement := `
		INSERT INTO sorena_counts (room_id, user_id, count)
		VALUES ($1, $2, $3)
		ON CONFLICT (room_id, user_id)
		DO UPDATE SET count = sorena_counts.count + EXCLUDED.count
	`
	_, err = h.db.Exec(sqlStatement, roomID, req.UserID, req.Count)
	if err != nil {
		log.Println("カウントの更新に失敗しました:", err)
		http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}