package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
    "github.com/joho/godotenv"
)

type MockAIResponse struct {
	InitialQuestion string `json:"initial_question"`
}

type RoomInfo struct {
	RoomID string `json:"room_id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

type Participant struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
}

type StartRoomResponse struct {
	InitialQuestion string       `json:"initial_question"`
	RoomInfo        RoomInfo     `json:"room_info"`
	Participants    []Participant `json:"participants"`
}

func TestStartRoomHandler(t *testing.T) {
    err := godotenv.Load()
    if err != nil {
        t.Fatal(".envファイルの読み込みに失敗しました")
    }

	dbPassword := os.Getenv("DB_PASSWORD")
    if dbPassword == "" {
        t.Fatal("環境変数 DB_PASSWORD が設定されていません")
    }

	// モックデータベース接続
	db, err := sql.Open("pgx", "user=postgres password=DB_PASSWORD dbname=elmo-db sslmode=disable")
	if err != nil {
		t.Fatalf("データベース接続エラー: %v", err)
	}
	defer db.Close()

	// モックリクエスト
	roomID := "r001"
	url := "/rooms/" + roomID + "/start"
	req := httptest.NewRequest(http.MethodPost, url, nil)

	// モックレスポンス
	w := httptest.NewRecorder()

	// RoomHandlerを使用してハンドラー関数を呼び出し
	handler := NewRoomHandler(db)
	handler.StartRoom(w, req)

	// レスポンスの検証
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("期待したステータスコード: %d, 実際のステータスコード: %d", http.StatusOK, res.StatusCode)
	}

	var response StartRoomResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("レスポンスのデコードエラー: %v", err)
	}

	// 初期質問の検証
	if response.InitialQuestion == "" {
		t.Error("初期質問が生成されていません")
	}

	// ルーム情報の検証
	if response.RoomInfo.RoomID != roomID {
		t.Errorf("期待したルームID: %s, 実際のルームID: %s", roomID, response.RoomInfo.RoomID)
	}
	if response.RoomInfo.Status != "inprogress" {
		t.Errorf("期待したステータス: inprogress, 実際のステータス: %s", response.RoomInfo.Status)
	}

	// 参加者情報の検証
	if len(response.Participants) == 0 {
		t.Error("参加者情報が空です")
	}
}
