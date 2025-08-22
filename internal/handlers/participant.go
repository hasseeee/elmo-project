package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin" // ★ Ginをインポート
	"github.com/shuto.sawaki/elmo-project/internal/models"
)

type ParticipantHandler struct {
	db *sql.DB
}

func NewParticipantHandler(db *sql.DB) *ParticipantHandler {
	return &ParticipantHandler{db: db}
}

// GET /participants
func (h *ParticipantHandler) GetParticipants(c *gin.Context) {
	// ★ クエリパラメータを取得します (例: /participants?room_id=r001)
	roomID := c.Query("room_id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "room_idは必須です"})
		return
	}
	// ... (DBから参加者リストを取得するロジック) ...
	var users []models.ParticipantUser
	// ...
	response := models.ParticipantsResponse{
		RoomID: roomID,
		Users:  users,
	}
	c.JSON(http.StatusOK, response)
}

// POST /participants
func (h *ParticipantHandler) AddParticipant(c *gin.Context) {
	var req models.ParticipantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "リクエストボディが不正です"})
		return
	}
	// ... (DBに参加者を追加するロジック) ...
	c.Status(http.StatusCreated)
}