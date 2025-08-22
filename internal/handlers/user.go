package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

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
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POSTメソッドのみサポートされています", http.StatusMethodNotAllowed)
		return
	}

	var newUser models.User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		// ... (以降のCreateUserのロジック)
	}
}