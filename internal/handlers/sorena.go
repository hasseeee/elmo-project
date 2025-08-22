package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/shuto.sawaki/elmo-project/internal/models"
)

func (h *RoomHandler) HandleSorena(w http.ResponseWriter, r *http.Request, roomID string) {
	var req models.SorenaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "リクエストボディが不正です", http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		http.Error(w, "user_idは必須です", http.StatusBadRequest)
		return
	}
	if req.Count <= 0 {
		http.Error(w, "countは1以上の値が必要です", http.StatusBadRequest)
		return
	}

	var exists bool
	err := h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM rooms WHERE id = $1)", roomID).Scan(&exists)
	if err != nil {
		log.Println("部屋の存在確認に失敗しました:", err)
		http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "指定された部屋は存在しません", http.StatusNotFound)
		return
	}

	err = h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", req.UserID).Scan(&exists)
	if err != nil {
		log.Println("ユーザーの存在確認に失敗しました:", err)
		http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "指定されたユーザーは存在しません", http.StatusNotFound)
		return
	}

	sqlStatement := `
		INSERT INTO sorena_counts (room_id, user_id, count)
		VALUES ($1, $2, $3)
		ON CONFLICT (room_id, user_id)
		DO UPDATE SET count = sorena_counts.count + EXCLUDED.count
	`
	_, err = h.db.Exec(sqlStatement, roomID, req.UserID, req.Count)
	if err != nil {
		log.Println("カウントの更新に失敗しました:", err)
		http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
