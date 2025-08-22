package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/matoous/go-nanoid/v2"
	"github.com/jackc/pgx/v5/pgconn"
	"strings"
	"errors"
)

type Room struct {
	ID          string    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type User struct {
	ID			string	`json:"id"`
	UserName	string `json:"user_name"`
}

var db *sql.DB

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		log.Fatal("DB_PASSWORD not set in .env file")
	}

connStr := fmt.Sprintf("user=postgres password=%s dbname=elmo-db sslmode=disable client_encoding=UTF8", password)

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
	http.HandleFunc("/users", usersHandler)
	http.HandleFunc("/users/all", getUsersHandler)
	log.Println("サーバー起動: http://localhost:8080")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("サーバーの起動に失敗しました:", err)
	}
}

func roomsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getRoomsHandler(w, r)
	case "POST":
		createRoomHandler(w, r)
	default:
		http.Error(w, "サポートされていないメソッドです", http.StatusMethodNotAllowed)
	}
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
	getUsersHandler(w, r)
	case "POST":	
		createUserHandler(w, r)
	default:
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

func createRoomHandler(w http.ResponseWriter, r *http.Request) {
	var newRoom Room
	err := json.NewDecoder(r.Body).Decode(&newRoom)
	if err != nil {
		http.Error(w, "リクエストボディが不正です", http.StatusBadRequest)
		return
	}
	
	if newRoom.Title == "" {
		http.Error(w, "タイトルは必須です", http.StatusBadRequest)
		return
	}

	newId, err := gonanoid.Generate("0123456789abcdefghijklmnopqrstuvwxyz", 6)
	if err != nil {
		log.Println("IDの生成に失敗しました:", err)
		http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
		return
	}
	newRoom.ID = newId

	sqlStatement := `INSERT INTO rooms (id, title, description) VALUES ($1, $2, $3)`
	_, err = db.Exec(sqlStatement, newRoom.ID, newRoom.Title, newRoom.Description)

	if err != nil {
		log.Println("データベースへのINSERTに失敗しました:", err)
		http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
		return
	}

	log.Printf("新しい部屋を作成しました: ID=%s, Title=%s", newRoom.ID, newRoom.Title)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newRoom)
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "リクエストボディが不正です", http.StatusBadRequest)
		return
	}

	fmt.Println(newUser)

	if newUser.UserName == "" {
		http.Error(w, "ユーザ名必須です", http.StatusBadRequest)
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
		_, err = db.Exec(sqlStatement, newUser.ID, newUser.UserName)

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

func getRoomByIdHandler(w http.ResponseWriter, r *http.Request) {

    id := strings.TrimPrefix(r.URL.Path, "/rooms/")

    var room Room
    sqlStatement := `SELECT id, title, description FROM rooms WHERE id = $1`
    err := db.QueryRow(sqlStatement, id).Scan(&room.ID, &room.Title, &room.Description)
    if err != nil {
        if err == sql.ErrNoRows {
            // もし行が見つからなかった場合 (sql.ErrNoRows)
            http.Error(w, "指定された部屋は見つかりません", http.StatusNotFound) // 404 Not Found
        } else {
            log.Println("データベースクエリの実行に失敗しました:", err)
            http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError) // 500 Internal Server Error
        }
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(room)
}

func getUsersHandler(w http.ResponseWriter, _ *http.Request) {
	rows, err := db.Query("SELECT id, user_name FROM users ORDER BY user_name ASC")
	if err != nil {
		log.Println("データベースクエリの実行に失敗しました:", err)
		http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.UserName)
		if err != nil {
			log.Println("データベースからのデータ読み取りに失敗しました:", err)
			http.Error(w, "サーバー内部エラーです", http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}