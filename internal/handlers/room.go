package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin" // ★ Ginをインポート
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
	// ... (GetRoomsのロジックはほぼ同じ) ...
	var rooms []models.Room
    // ...
	// ★ レスポンスの書き方がシンプルになります
	c.JSON(http.StatusOK, rooms)
}

// POST /rooms
func (h *RoomHandler) CreateRoom(c *gin.Context) {
	var newRoom models.Room
	// ★ リクエストボディのJSONを構造体にバインドします
	if err := c.ShouldBindJSON(&newRoom); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "リクエストボディが不正です"})
		return
	}
	// ... (ID生成やDBへのINSERT処理) ...
	c.JSON(http.StatusCreated, newRoom)
}

// GET /rooms/:id
func (h *RoomHandler) GetRoomByID(c *gin.Context) {
	// ★ URLパラメータ(:id)を取得します
	id := c.Param("id")
	var room models.Room
	// ... (DBから部屋情報を取得するロジック) ...
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
	// ... (結論をDBに保存するロジック) ...
	var room models.Room // 更新後の部屋情報を取得
	// ...
	c.JSON(http.StatusOK, room)
}

// GET /rooms/:id/start
func (h *RoomHandler) StartRoom(c *gin.Context) {
	roomID := c.Param("id")
	// ... (部屋のステータスチェック、AIからの質問生成、DB更新のロジック) ...
	var participants []models.ParticipantUser
	// ... (参加者リスト取得のロジック) ...
	response := models.StartRoomResponse{
		// ...
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
	// ... (sorenaカウントをDBに保存するロジック) ...
	c.Status(http.StatusNoContent) // ボディなしの成功レスポンス
}