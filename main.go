package main

import (
	"database/sql"
	"encoding/json"
	"fmt" // fmtパッケージを追加
	"log"
	"net/http"
	"os" // osパッケージを追加

	// godotenvパッケージをインポート
	"github.com/joho/godotenv"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// (Room構造体は変更なし)
type Room struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

var db *sql.DB

func main() {
	// .envファイルを読み込む
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 環境変数からパスワードを取得
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		log.Fatal("DB_PASSWORD not set in .env file")
	}

	// 取得したパスワードを使って接続情報を作成
	connStr := fmt.Sprintf("user=postgres password=%s dbname=chat_app_db sslmode=disable", password)

	// データベースへの接続 (以降は変更なし)
	db, err = sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal("データベースへの接続に失敗しました:", err)
	}
	defer db.Close()
	
	// ... (以降のmain関数のコードは変更なし) ...
	err = db.Ping()
	if err != nil {
		log.Fatal("データベースへの疎通確認に失敗しました:", err)
	}
	log.Println("データベースへの接続に成功しました。")
	http.HandleFunc("/rooms", getRoomsHandler)
	log.Println("サーバー起動: http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("サーバーの起動に失敗しました:", err)
	}
}

// (getRoomsHandler関数は変更なし)
func getRoomsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, title, description FROM rooms ORDER BY id ASC")
	if err != nil {
		log.Println("データベースクエリの実行に失敗しました:", err)
		http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var rooms []Room

	for rows.Next() {
		var room Room
		err := rows.Scan(&room.ID, &room.Title, &room.Description)
		if err != nil {
			log.Println("データベースからのデータ読み取りに失敗しました:", err)
			http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
			return
		}
		rooms = append(rooms, room)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(rooms)
	if err != nil {
		log.Println("JSONへの変換に失敗しました:", err)
	}
}