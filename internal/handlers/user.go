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

func (h *UserHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		h.CreateUser(w, r)
	default:
		http.Error(w, "サポートされていないメソッドです", http.StatusMethodNotAllowed)
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser models.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "リクエストボディが不正です", http.StatusBadRequest)
		return
	}

	if newUser.UserName == "" {
		http.Error(w, "ユーザ名は必須です", http.StatusBadRequest)
		return
	}

	const maxRetries = 10
	for i := 0; i < maxRetries; i++ {
		newId, err := gonanoid.Generate("0123456789abcdefghijklmnopqrstuvwxyz", 8)
		if err != nil {
			log.Println("IDの生成に失敗しました:", err)
			http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
			return
		}
		newUser.ID = newId

		sqlStatement := `INSERT INTO users (id, user_name) VALUES ($1, $2)`
		_, err = h.db.Exec(sqlStatement, newUser.ID, newUser.UserName)

		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				log.Printf("User IDが重複しました。再試行します... (試行 %d回目)", i+1)
				continue
			}

			log.Println("データベースへのINSERTに失敗しました:", err)
			http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
			return
		}

		log.Printf("新しいユーザを作成しました: ID=%s, UserName=%s", newUser.ID, newUser.UserName)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newUser)
		return
	}

	log.Println("User IDの生成に最大回数失敗しました。")
	http.Error(w, "サーバー内部で問題が発生しました。", http.StatusInternalServerError)
}
