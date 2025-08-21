package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	// PostgreSQLのドライバーを読み込む
	// _ は、このパッケージの関数を直接は使わないが、内部的に必要だというGoのおまじない
	_ "github.com/jackc/pgx/v5/stdlib" 
)

// Room 構造体は変更なし
type Room struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// データベース接続を保持するための変数
var db *sql.DB

func main() {
	// 1. データベースへの接続情報を作成する（★必ず自分の設定に書き換える★）
	connStr := "user=postgres password=あなたのパスワード dbname=chat_app_db sslmode=disable"

	var err error
	// 2. PostgreSQLに接続する
	// "pgx"は使用するドライバー名
	db, err = sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal("データベースへの接続に失敗しました:", err)
	}
	// main関数が終わる時に、必ずデータベース接続を閉じる
	defer db.Close()

	// 3. データベースへの接続確認
	err = db.Ping()
	if err != nil {
		log.Fatal("データベースへの疎通確認に失敗しました:", err)
	}

	log.Println("データベースへの接続に成功しました。")

	// 4. APIサーバーを起動する
	http.HandleFunc("/rooms", getRoomsHandler)

	log.Println("サーバー起動: http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("サーバーの起動に失敗しました:", err)
	}
}

// 部屋一覧を取得するハンドラ関数
func getRoomsHandler(w http.ResponseWriter, r *http.Request) {
	// 1. データベースから全ルームの情報を取得するSQLクエリを実行
	rows, err := db.Query("SELECT id, title, description FROM rooms ORDER BY id ASC")
	if err != nil {
		// もしクエリ実行でエラーが起きたら
		log.Println("データベースクエリの実行に失敗しました:", err)
		// サーバー内部でエラーが起きたことを示す500エラーを返す
		http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
		return
	}
	// この関数が終わる時に、必ずrowsを閉じる
	defer rows.Close()

	// 2. 取得したデータを格納するためのGoのスライスを用意
	var rooms []Room

	// 3. 1行ずつデータを取り出す
	for rows.Next() {
		var room Room
		// 取り出したデータをRoom構造体の各フィールドに割り当てる
		err := rows.Scan(&room.ID, &room.Title, &room.Description)
		if err != nil {
			log.Println("データベースからのデータ読み取りに失敗しました:", err)
			http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
			return
		}
		// 割り当てたRoomをスライスに追加
		rooms = append(rooms, room)
	}

	// 4. 最終的な結果をJSON形式でクライアントに返す
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(rooms)
	if err != nil {
		log.Println("JSONへの変換に失敗しました:", err)
	}
}