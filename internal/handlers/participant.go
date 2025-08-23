package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shuto.sawaki/elmo-project/internal/models"
)

type ParticipantHandler struct {
	db *sql.DB
}

func NewParticipantHandler(db *sql.DB) *ParticipantHandler {
	return &ParticipantHandler{db: db}
}

// GetParticipants godoc
// @Summary      参加者一覧を取得
// @Description  指定された会議室の参加者一覧を取得します
// @Tags         participants
// @Accept       json
// @Produce      json
// @Param        room_id  query     string  true  "会議室ID"
// @Success      200      {object}  models.ParticipantsResponse
// @Failure      400      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /participants [get]
func (h *ParticipantHandler) GetParticipants(c *gin.Context) {
	roomID := c.Query("room_id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "room_idは必須です"})
		return
	}

	rows, err := h.db.Query("SELECT u.id, u.user_name FROM participants p JOIN users u ON p.user_id = u.id WHERE p.room_id = $1", roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバー内部エラーです"})
		return
	}
	defer rows.Close()

	var users []models.ParticipantUser
	for rows.Next() {
		var user models.ParticipantUser
		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバー内部エラーです"})
			return
		}
		users = append(users, user)
	}

	response := models.ParticipantsResponse{
		RoomID: roomID,
		Users:  users,
	}
	c.JSON(http.StatusOK, response)
}

// AddParticipant godoc
// @Summary      参加者を追加
// @Description  指定された会議室に参加者を追加します
// @Tags         participants
// @Accept       json
// @Produce      json
// @Param        participant  body      models.ParticipantRequest  true  "参加者情報"
// @Success      201          "Created"
// @Failure      400          {object}  map[string]interface{}
// @Failure      500          {object}  map[string]interface{}
// @Router       /participants [post]
func (h *ParticipantHandler) AddParticipant(c *gin.Context) {
	var req models.ParticipantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "リクエストボディが不正です"})
		return
	}

	if req.RoomID == "" || req.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "room_idとuser_idは必須です"})
		return
	}

	_, err := h.db.Exec("INSERT INTO participants (room_id, user_id) VALUES ($1, $2)", req.RoomID, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "参加者の追加に失敗しました"})
		return
	}
	c.Status(http.StatusCreated)
}