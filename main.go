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
	"github.com/matoous/go-nanoid/v2"
	"strings"
)

type Room struct {
	ID          string    `json:"id"`
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
    http.HandleFunc("/rooms/", getRoomByIdHandler) 
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
	// 6文字のIDをGoで生成する
	newId, err := gonanoid.Generate("0123456789abcdefghijklmnopqrstuvwxyz", 6)
	if err != nil {
		log.Println("IDの生成に失敗しました:", err)
		http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
		return
	}
	newRoom.ID = newId // 生成したIDをセット

	// ★変更点：生成したIDを含めてINSERTする
	sqlStatement := `INSERT INTO rooms (id, title, description) VALUES ($1, $2, $3)`
	_, err = db.Exec(sqlStatement, newRoom.ID, newRoom.Title, newRoom.Description)

	if err != nil {
		log.Println("データベースへのINSERTに失敗しました:", err)
		http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
		return
	}

	log.Printf("新しい部屋を作成しました: ID=%s, Title=%s", newRoom.ID, newRoom.Title)

	// 3. 成功したことをクライアントに伝える
	// 201 Createdステータスを返し、作成されたリソースをJSONで返す
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newRoom)
}

// IDを指定して特定の部屋を1件取得するハンドラ
func getRoomByIdHandler(w http.ResponseWriter, r *http.Request) {
    // 1. URLからID部分を文字列として抜き出す
    // 例：/rooms/abcdef -> "abcdef" を取得
    id := strings.TrimPrefix(r.URL.Path, "/rooms/")

    // 2. データベースに問い合わせて、指定されたIDの部屋を取得
    var room Room
    sqlStatement := `SELECT id, title, description FROM rooms WHERE id = $1`
    // QueryRowは、結果が1行だけのクエリに使うと便利
    err := db.QueryRow(sqlStatement, id).Scan(&room.ID, &room.Title, &room.Description)
    if err != nil {
        // データベースからのエラーを判定
        if err == sql.ErrNoRows {
            // もし行が見つからなかった場合 (sql.ErrNoRows)
            http.Error(w, "指定された部屋は見つかりません", http.StatusNotFound) // 404 Not Found
        } else {
            // その他のデータベースエラー
            log.Println("データベースクエリの実行に失敗しました:", err)
            http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError) // 500 Internal Server Error
        }
        return
    }

    // 3. 見つかった部屋の情報をJSONで返す
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(room)
}