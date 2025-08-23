package handlers

import (
	"database/sql"
	"net/http"
	"log"
	"sync"
	"time"
	"errors"

	"github.com/gin-gonic/gin"
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
func (h *RoomHandler) GetRooms(c *gin.Context) {
	rows, err := h.db.Query("SELECT id, title, description FROM rooms ORDER BY id ASC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバー内部エラーです"})
		return
	}
	defer rows.Close()

	var rooms []models.Room
	for rows.Next() {
		var room models.Room
		if err := rows.Scan(&room.ID, &room.Title, &room.Description); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバー内部エラーです"})
			return
		}
		rooms = append(rooms, room)
	}
	c.JSON(http.StatusOK, rooms)
}

// POST /rooms
func (h *RoomHandler) CreateRoom(c *gin.Context) {
	var newRoom models.Room
	if err := c.ShouldBindJSON(&newRoom); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "リクエストボディが不正です"})
		return
	}

	if newRoom.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "タイトルは必須です"})
		return
	}

	newId, err := gonanoid.Generate("0123456789abcdefghijklmnopqrstuvwxyz", 6)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバー内部エラーです"})
		return
	}
	newRoom.ID = newId

	sqlStatement := `INSERT INTO rooms (id, title, description) VALUES ($1, $2, $3)`
	_, err = h.db.Exec(sqlStatement, newRoom.ID, newRoom.Title, newRoom.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバー内部エラーです"})
		return
	}
	c.JSON(http.StatusCreated, newRoom)
}

// GET /rooms/:id
func (h *RoomHandler) GetRoomByID(c *gin.Context) {
	id := c.Param("id")
	var room models.Room
	sqlStatement := `SELECT id, title, description, conclusion, status, initial_question FROM rooms WHERE id = $1`
	err := h.db.QueryRow(sqlStatement, id).Scan(&room.ID, &room.Title, &room.Description, &room.Conclusion, &room.Status, &room.InitialQuestion)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "指定された部屋は見つかりません"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバー内部エラーです"})
		}
		return
	}
	c.JSON(http.StatusOK, room)
}

// POST /rooms/:id/conclusion
func (h *RoomHandler) SaveConclusion(c *gin.Context) {
	roomID := c.Param("id")
	var req models.ConclusionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "リクエストボディが不正です"})
		return
	}
	if req.Conclusion == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "結論は必須です"})
		return
	}
	_, err := h.db.Exec("UPDATE rooms SET conclusion = $1, status = 'concluded' WHERE id = $2", req.Conclusion, roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバー内部エラーです"})
		return
	}
	var room models.Room
	err = h.db.QueryRow("SELECT id, title, description, conclusion, status FROM rooms WHERE id = $1", roomID).
		Scan(&room.ID, &room.Title, &room.Description, &room.Conclusion, &room.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバー内部エラーです"})
		return
	}
	c.JSON(http.StatusOK, room)
}

// GET /rooms/:id/start
func (h *RoomHandler) StartRoom(c *gin.Context) {
	roomID := c.Param("id")

	var room models.Room
	err := h.db.QueryRow("SELECT id, title, description, status FROM rooms WHERE id = $1", roomID).Scan(&room.ID, &room.Title, &room.Description, &room.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバー内部エラーです"})
		return
	}

	if room.Status != "not started" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "部屋は既に開始されています"})
		return
	}

	initialQuestion, err := h.aiGenerator.GenerateInitialQuestion(c.Request.Context(), room.Title, room.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI API呼び出しエラー"})
		return
	}

	_, err = h.db.Exec("UPDATE rooms SET status = $1, initial_question = $2 WHERE id = $3", "inprogress", initialQuestion, roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバー内部エラーです"})
		return
	}

	rows, err := h.db.Query(`SELECT u.id, u.user_name FROM participants p JOIN users u ON p.user_id = u.id WHERE p.room_id = $1`, roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバー内部エラーです"})
		return
	}
	defer rows.Close()

	var participants []models.ParticipantUser
	for rows.Next() {
		var p models.ParticipantUser
		if err := rows.Scan(&p.ID, &p.Name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバー内部エラーです"})
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
	c.JSON(http.StatusOK, response)
}

// POST /rooms/:id/sorena
func (h *RoomHandler) HandleSorena(c *gin.Context) {
	roomID := c.Param("id")
	var req models.SorenaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "リクエストボディが不正です"})
		return
	}

	// ... (元のsorena.goにあったバリデーションを追加しても良い) ...

	sqlStatement := `
		INSERT INTO sorena_counts (room_id, user_id, count)
		VALUES ($1, $2, $3)
		ON CONFLICT (room_id, user_id)
		DO UPDATE SET count = sorena_counts.count + EXCLUDED.count
	`
	_, err := h.db.Exec(sqlStatement, roomID, req.UserID, req.Count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバー内部エラーです"})
		return
	}
	c.Status(http.StatusNoContent)
}

// POST /rooms/:id/summary
func (h *RoomHandler) CreateSummary(c *gin.Context) {
	// 1. URLから部屋のIDを取得
	roomID := c.Param("id")

	// 2. リクエストのJSONデータをGoの構造体に変換
	var req models.SummaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "リクエストボディが不正です"})
		return
	}

	// ログが空の場合は何もしない
	if len(req.Logs) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	// 3. AIに要約を依頼
	summary, err := h.aiGenerator.SummarizeLogs(c.Request.Context(), req.Logs)
	if err != nil {
		// ここではエラーをログに出力するだけにして、クライアントにはエラーを返さないことも考えられます。
		// 定期実行のバックグラウンド処理的な側面が強いため。今回はサーバーエラーとして返します。
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AIによる要約に失敗しました"})
		return
	}

	// 4. 要約結果をDBに保存
	logID, err := gonanoid.New() // 要約ログの新しいIDを生成
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "IDの生成に失敗しました"})
		return
	}

	sqlStatement := `
		INSERT INTO chat_logs (id, room_id, message, is_summary)
		VALUES ($1, $2, $3, TRUE)
	`
	_, err = h.db.Exec(sqlStatement, logID, roomID, summary)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "データベースへの保存に失敗しました"})
		return
	}

	// 5. 成功したが返すコンテンツはない、というステータスを返す
	c.Status(http.StatusNoContent)
}

// PUT /rooms/:id/status
func (h *RoomHandler) UpdateRoomStatus(c *gin.Context) {
	// URLからidを取得
	roomID := c.Param("id")

	var req models.UpdateRoomStatusRequest
	// リクエストボディのJSONを構造体にバインド（割り当て）
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// バリデーション：statusが"done"であることのみを許可
	if req.Status != "done" {
		c.JSON(http.StatusBadRequest, gin.H{"error": `status must be "done"`})
		return
	}

	// データベースを更新するSQL
	sqlStatement := `UPDATE rooms SET status = $1 WHERE id = $2`
	result, err := h.db.ExecContext(c.Request.Context(), sqlStatement, req.Status, roomID)
	if err != nil {
		log.Printf("failed to update room status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// 更新対象の行が存在したか確認
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("failed to get rows affected: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}

	// 成功時は 204 No Content を返す
	c.Status(http.StatusNoContent)
}

// GET /rooms/:id/result
func (h *RoomHandler) GetRoomResult(c *gin.Context) {
	roomID := c.Param("id")
	ctx := c.Request.Context()

	var wg sync.WaitGroup
	var roomInfo models.ResultRoomInfo
	var sorenaSummary models.SorenaSummary
	var chatLogs []models.ChatLog // ★ LogEntry から ChatLog に変更
	var errRoom, errSorena, errLogs error

	wg.Add(3)

	// Goroutine 1: 部屋情報を取得
	go func() {
		defer wg.Done()
		var title string
		err := h.db.QueryRowContext(ctx, "SELECT title FROM rooms WHERE id = $1", roomID).Scan(&title)
		if err != nil {
			errRoom = err
			return
		}
		roomInfo = models.ResultRoomInfo{RoomID: roomID, Title: title}
	}()

	// Goroutine 2: 「それな」の集計
	go func() {
		defer wg.Done()
		query := `
            SELECT u.id, u.user_name, COUNT(sc.id) as count
            FROM sorena_counts sc
            JOIN users u ON sc.user_id = u.id
            WHERE sc.room_id = $1
            GROUP BY u.id, u.user_name
            ORDER BY count DESC`
		rows, err := h.db.QueryContext(ctx, query, roomID)
		if err != nil {
			errSorena = err
			return
		}
		defer rows.Close()

		var participants []models.SorenaParticipant
		totalCount := 0
		for rows.Next() {
			var p models.SorenaParticipant
			if err := rows.Scan(&p.UserID, &p.UserName, &p.Count); err != nil {
				errSorena = err
				return
			}
			participants = append(participants, p)
			totalCount += p.Count
		}
		sorenaSummary = models.SorenaSummary{
			TotalCount:   totalCount,
			Participants: participants,
		}
	}()

	// Goroutine 3: チャットログを取得
	go func() {
		defer wg.Done()
		query := `
			SELECT id, user_id, content, is_summary, created_at
			FROM chat_logs
			WHERE room_id = $1
			ORDER BY created_at ASC`
		rows, err := h.db.QueryContext(ctx, query, roomID)
		if err != nil {
			errLogs = err
			return
		}
		defer rows.Close()

		for rows.Next() {
			var log models.ChatLog // ★ LogEntry から ChatLog に変更
			var userID sql.NullString
			// ★ ScanするフィールドをChatLogのフィールド名に合わせる
			if err := rows.Scan(&log.ID, &userID, &log.Message, &log.IsSummary, &log.Timestamp); err != nil {
				errLogs = err
				return
			}
			if userID.Valid {
				log.UserID = &userID.String
			}
			chatLogs = append(chatLogs, log)
		}
	}()

	wg.Wait()

	if errRoom != nil || errSorena != nil || errLogs != nil {
		log.Printf("Error fetching room result: roomErr=%v, sorenaErr=%v, logErr=%v", errRoom, errSorena, errLogs)
		if errors.Is(errRoom, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch room result data"})
		return
	}

	response := models.RoomResultResponse{
		RoomInfo:      roomInfo,
		SorenaSummary: sorenaSummary,
		ChatLogs:      chatLogs,
	}

	c.JSON(http.StatusOK, response)
}