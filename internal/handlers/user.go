package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
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
		newId, err := gonanoid.Generate("0123456789abcdefghijklmnopqrstuvwxyz", 8)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバー内部エラーです"})
			return
		}
		newUser.ID = newId

		sqlStatement := `INSERT INTO users (id, user_name) VALUES ($1, $2)`
		_, err = h.db.Exec(sqlStatement, newUser.ID, newUser.UserName)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				continue // IDが重複した場合はループを継続して再試行
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバー内部エラーです"})
			return
		}

		c.JSON(http.StatusCreated, newUser)
		return // 成功したらループを抜ける
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバー内部で問題が発生しました。"})
}