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
	connStr := fmt.Sprintf("user=postgres password=%s dbname=elmo-db sslmode=disable", password)

	// データベースへの接続 (以降は変更なし)
	db, err = sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal("データベースへの接続に失敗しました:", err)
	}
	defer db.Close()
	
	err = db.Ping()
	if err != nil {
		log.Fatal("データベースへの疎通確認に失敗しました:", err)
	}
	log.Println("データベースへの接続に成功しました。")

	http.HandleFunc("/rooms", roomsHandler)
	log.Println("サーバー起動: http://localhost:8080")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("サーバーの起動に失敗しました:", err)
	}
}

// roomsHandlerは、リクエストの種類(GET/POST)に応じて処理を振り分ける
func roomsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getRoomsHandler(w, r)
	case "POST":
		createRoomHandler(w, r)
	default:
		// GETとPOST以外のメソッドが来たら、エラーを返す
		http.Error(w, "サポートされていないメソッドです", http.StatusMethodNotAllowed)
	}
}

func getRoomsHandler(w http.ResponseWriter, _ *http.Request) {
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

// 新しい部屋を作成するハンドラ
func createRoomHandler(w http.ResponseWriter, r *http.Request) {
	var newRoom Room
	// 1. リクエストのJSONボディをデコードして、newRoomにマッピング
	err := json.NewDecoder(r.Body).Decode(&newRoom)
	if err != nil {
		http.Error(w, "リクエストボディが不正です", http.StatusBadRequest)
		return
	}
	
	// タイトルが空の場合はエラー
	if newRoom.Title == "" {
		http.Error(w, "タイトルは必須です", http.StatusBadRequest)
		return
	}

	// 2. データベースに新しい部屋をINSERTする
	// QueryRowを使って、INSERTした行のIDを取得する
	sqlStatement := `INSERT INTO rooms (title, description) VALUES ($1, $2) RETURNING id`
	err = db.QueryRow(sqlStatement, newRoom.Title, newRoom.Description).Scan(&newRoom.ID)
	if err != nil {
		log.Println("データベースへのINSERTに失敗しました:", err)
		http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
		return
	}

	log.Printf("新しい部屋を作成しました: ID=%d, Title=%s", newRoom.ID, newRoom.Title)

	// 3. 成功したことをクライアントに伝える
	// 201 Createdステータスを返し、作成されたリソースをJSONで返す
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newRoom)
}
