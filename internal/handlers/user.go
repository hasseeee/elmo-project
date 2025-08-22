package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin" // ★ Ginをインポート
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/matoous/go-nanoid/v2"
	"github.com/shuto.sawaki/elmo-project/internal/models"
)

type UserHandler struct {
	db *sql.DB
}

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{db: db}
}

// POST /users
func (h *UserHandler) CreateUser(c *gin.Context) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "リクエストボディが不正です"})
		return
	}

	if newUser.UserName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ユーザ名必須です"})
		return
	}

	const maxRetries = 10
	for i := 0; i < maxRetries; i++ {
		// ... (ID生成とDBへのINSERTロジックは同じ) ...
		// 成功した場合
		c.JSON(http.StatusCreated, newUser)
		return
	}

	// 失敗した場合
	c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバー内部で問題が発生しました。"})
}