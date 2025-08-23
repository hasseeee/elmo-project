package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin" // ★ Ginをインポート
	"github.com/shuto.sawaki/elmo-project/internal/ai"
	"github.com/shuto.sawaki/elmo-project/internal/db"
	"github.com/shuto.sawaki/elmo-project/internal/handlers"
)

func main() {
	database, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	ctx := context.Background()
	aiGenerator, err := ai.NewGeminiAIGenerator(ctx)
	if err != nil {
		log.Fatalf("AIジェネレータの初期化に失敗しました: %v", err)
	}

	// 各ハンドラーを初期化
	roomHandler := handlers.NewRoomHandler(database, aiGenerator)
	userHandler := handlers.NewUserHandler(database)
	participantHandler := handlers.NewParticipantHandler(database)

	// ★ Ginのルーターを初期化
	// gin.Default()は、ロガーやリカバリーといった便利なミドルウェアが最初から組み込まれています。
	router := gin.Default()

	// --- ルーティングの設定 ---
	// GETとPOSTのようにHTTPメソッドごとに明確にルートを定義できます。
	// これにより、ハンドラー内のswitch文が不要になります。
	router.GET("/rooms", roomHandler.GetRooms)
	router.POST("/rooms", roomHandler.CreateRoom)

	// URL内の可変部分を :id のようにコロンで指定できます。
	// これを「URLパラメータ」と呼びます。
	router.GET("/rooms/:id", roomHandler.GetRoomByID)
	router.POST("/rooms/:id/start", roomHandler.StartRoom)
	router.POST("/rooms/:id/conclusion", roomHandler.SaveConclusion)
	router.POST("/rooms/:id/sorena", roomHandler.HandleSorena)
	router.POST("/rooms/:id/summary", roomHandler.CreateSummary)

	router.POST("/users", userHandler.CreateUser)

	router.GET("/participants", participantHandler.GetParticipants)
	router.POST("/participants", participantHandler.AddParticipant)

	log.Println("サーバー起動: http://localhost:8080")
	// ★ Ginのルーターでサーバーを起動
	router.Run(":8080")
}