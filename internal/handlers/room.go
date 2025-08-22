package handlers

import (
	"database/sql"
	"net/http"

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