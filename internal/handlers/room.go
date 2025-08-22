package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/matoous/go-nanoid/v2"
	"github.com/shuto.sawaki/elmo-project/internal/models"
)

type RoomHandler struct {
	db *sql.DB
}

func NewRoomHandler(db *sql.DB) *RoomHandler {
	return &RoomHandler{db: db}
}

func (h *RoomHandler) HandleRooms(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.GetRooms(w, r)
	case "POST":
		h.CreateRoom(w, r)
	default:
		http.Error(w, "サポートされていないメソッドです", http.StatusMethodNotAllowed)
	}
}

func (h *RoomHandler) GetRooms(w http.ResponseWriter, _ *http.Request) {
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
	err = json.NewEncoder(w).Encode(rooms)
	if err != nil {
		log.Println("JSONへの変換に失敗しました:", err)
	}
}

func (h *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var newRoom models.Room
	err := json.NewDecoder(r.Body).Decode(&newRoom)
	if err != nil {
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

func (h *RoomHandler) GetRoomByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/rooms/")

	if strings.HasSuffix(id, "/sorena") {
		if r.Method == "POST" {
			roomID := strings.TrimSuffix(id, "/sorena")
			h.HandleSorena(w, r, roomID)
			return
		}
		http.Error(w, "サポートされていないメソッドです", http.StatusMethodNotAllowed)
		return
	}

	if strings.HasSuffix(id, "/conclusion") {
		if r.Method == "POST" {
			roomID := strings.TrimSuffix(id, "/conclusion")
			h.SaveConclusion(w, r, roomID)
			return
		}
		http.Error(w, "サポートされていないメソッドです", http.StatusMethodNotAllowed)
		return
	}

	var room models.Room
	sqlStatement := `SELECT id, title, description, conclusion FROM rooms WHERE id = $1`
	err := h.db.QueryRow(sqlStatement, id).Scan(&room.ID, &room.Title, &room.Description, &room.Conclusion)
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

func (h *RoomHandler) SaveConclusion(w http.ResponseWriter, r *http.Request, roomID string) {
	// リクエストボディを読み取り
	var req models.ConclusionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "リクエストボディが不正です", http.StatusBadRequest)
		return
	}

	// バリデーション
	if req.Conclusion == "" {
		http.Error(w, "結論は必須です", http.StatusBadRequest)
		return
	}

	// 部屋の存在確認
	var exists bool
	err := h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM rooms WHERE id = $1)", roomID).Scan(&exists)
	if err != nil {
		log.Println("部屋の存在確認に失敗しました:", err)
		http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "指定された部屋は見つかりません", http.StatusNotFound)
		return
	}

	// 結論を保存
	_, err = h.db.Exec("UPDATE rooms SET conclusion = $1 WHERE id = $2", req.Conclusion, roomID)
	if err != nil {
		log.Println("結論の保存に失敗しました:", err)
		http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
		return
	}

	// 更新後の部屋の情報を取得
	var room models.Room
	err = h.db.QueryRow("SELECT id, title, description, conclusion FROM rooms WHERE id = $1", roomID).
		Scan(&room.ID, &room.Title, &room.Description, &room.Conclusion)
	if err != nil {
		log.Println("更新後の部屋情報の取得に失敗しました:", err)
		http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
		return
	}

	// レスポンスを返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(room)
}
