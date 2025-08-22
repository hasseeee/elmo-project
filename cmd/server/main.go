package main

import (
	"log"
	"net/http"

	"github.com/shuto.sawaki/elmo-project/internal/db"
	"github.com/shuto.sawaki/elmo-project/internal/handlers"
)

func main() {
	// データベース接続の初期化
	database, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	// ハンドラーの初期化
	roomHandler := handlers.NewRoomHandler(database)
	userHandler := handlers.NewUserHandler(database)
	participantHandler := handlers.NewParticipantHandler(database)

	// ルーティングの設定
	http.HandleFunc("/rooms", roomHandler.HandleRooms)
	http.HandleFunc("/rooms/", roomHandler.GetRoomByID)
	http.HandleFunc("/users", userHandler.HandleUsers)
	http.HandleFunc("/participants", participantHandler.HandleParticipants)

	// サーバーの起動
	log.Println("サーバー起動: http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("サーバーの起動に失敗しました:", err)
	}
}
