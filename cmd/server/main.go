package main

import (
	"context" // contextパッケージをインポート
	"log"
	"net/http"

	"github.com/shuto.sawaki/elmo-project/internal/ai" // aiパッケージをインポート
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

	// --- ここからが修正箇所 ---
	// 本物のAIジェネレータを初期化
	ctx := context.Background()
	aiGenerator, err := ai.NewGeminiAIGenerator(ctx)
	if err != nil {
		log.Fatalf("AIジェネレータの初期化に失敗しました: %v", err)
	}
	// --- ここまで ---

	// ハンドラーの初期化
	// ★ 修正：NewRoomHandlerにaiGeneratorを2つ目の引数として渡す
	roomHandler := handlers.NewRoomHandler(database, aiGenerator)
	userHandler := handlers.NewUserHandler(database)
	participantHandler := handlers.NewParticipantHandler(database)

	// ルーティングの設定
	http.HandleFunc("/rooms", roomHandler.HandleRooms)
	http.HandleFunc("/rooms/", roomHandler.HandleRoomRequests)
	http.HandleFunc("/users", userHandler.HandleUsers)
	http.HandleFunc("/participants", participantHandler.HandleParticipants)

	// サーバーの起動
	log.Println("サーバー起動: http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("サーバーの起動に失敗しました:", err)
	}
}