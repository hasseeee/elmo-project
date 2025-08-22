package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/shuto.sawaki/elmo-project/internal/models"
	"google.golang.org/genai"
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
	log.Println("GetRooms: リクエストを受信しました")
	rows, err := h.db.Query("SELECT id, title, description FROM rooms ORDER BY id ASC")
	if err != nil {
		log.Println("データベースクエリの実行に失敗しました:", err)
		http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
		return
	}
	log.Println("GetRooms: クエリが成功しました")
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

func (h *RoomHandler) HandleRoomRequests(w http.ResponseWriter, r *http.Request) {
	log.Println("HandleRoomRequests: パス =", r.URL.Path)
	path := strings.TrimPrefix(r.URL.Path, "/rooms/")

	if strings.HasSuffix(path, "/start") {
		if r.Method == "GET" {
			h.StartRoom(w, r)
			return
		}
		http.Error(w, "サポートされていないメソッドです", http.StatusMethodNotAllowed)
		return
	}

	// それ以外のリクエストは既存のGetRoomByIDの処理を行う
	h.GetRoomByID(w, r)
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

func (h *RoomHandler) StartRoom(w http.ResponseWriter, r *http.Request) {
	log.Println("StartRoom: リクエストを受信しました")
	roomID := strings.TrimPrefix(r.URL.Path, "/rooms/")
	roomID = strings.TrimSuffix(roomID, "/start")

	// 部屋の存在と状態を確認
	var room models.Room
	sqlStatement := `SELECT id, title, description, status FROM rooms WHERE id = $1`
	err := h.db.QueryRow(sqlStatement, roomID).Scan(&room.ID, &room.Title, &room.Description, &room.Status)
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

	// initialQuestion を定義
	var initialQuestion string

	// GeminiAPI クライアントの初期化
	apiKey := os.Getenv("GOOGLE_API_KEY")
	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Println("Gemini クライアントの初期化に失敗しました:", err)
		http.Error(w, "AI API 初期化エラー", http.StatusInternalServerError)
		return
	}

	// AI API を呼び出して「はじめの問いかけ」を生成
	prompt := fmt.Sprintf("ルームタイトル: %s\n説明: %s\nこのルームに適した最初の問いかけを生成してください。", room.Title, room.Description)
	response, err := client.Models.GenerateContent(context.Background(), "gemini-1.5-flash", []*genai.Content{
		{Parts: []*genai.Part{{Text: prompt}}},
	}, nil)
	if err != nil {
		log.Println("GeminiAPI の呼び出しに失敗しました:", err)
		http.Error(w, "AI API 呼び出しエラー", http.StatusInternalServerError)
		return
	}

	if len(response.Candidates) == 0 {
		log.Println("GeminiAPI からの応答が空でした")
		http.Error(w, "AI API 応答エラー", http.StatusInternalServerError)
		return
	}

	if len(response.Candidates[0].Content.Parts) == 0 {
		log.Println("GeminiAPI からの応答に内容が含まれていませんでした")
		http.Error(w, "AI API 応答エラー", http.StatusInternalServerError)
		return
	}

	initialQuestion = response.Candidates[0].Content.Parts[0].Text // 最初の候補を使用
	log.Printf("生成された初期質問: %s", initialQuestion)

	// 部屋の状態と初期質問を更新
	updateStatement := `UPDATE rooms SET status = $1, initial_question = $2 WHERE id = $3`
	_, err = h.db.Exec(updateStatement, "inprogress", initialQuestion, roomID)
	if err != nil {
		log.Println("部屋の更新に失敗しました:", err)
		http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
		return
	}

	// 修正: ParticipantUser を使用して参加者リストを取得
	participants := []models.ParticipantUser{}
	rows, err := h.db.Query("SELECT id, name FROM participants WHERE room_id = $1", roomID)
	if err != nil {
		log.Println("参加者リストの取得に失敗しました:", err)
		http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var participant models.ParticipantUser
		if err := rows.Scan(&participant.ID, &participant.Name); err != nil {
			log.Println("参加者データの読み取りに失敗しました:", err)
			http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
			return
		}
		participants = append(participants, participant)
	}

	// 修正: 再代入時に `=` を使用
	responseData := map[string]interface{}{
		"initial_question": initialQuestion,
		"room_info": map[string]interface{}{
			"room_id": room.ID,
			"title":   room.Title,
			"status":  "inprogress",
		},
		"participants": participants,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}
