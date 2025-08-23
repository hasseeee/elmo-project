package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin" // ★ Ginをインポート
	"github.com/shuto.sawaki/elmo-project/internal/ai"
	"github.com/shuto.sawaki/elmo-project/internal/db"
	"github.com/shuto.sawaki/elmo-project/internal/handlers"
	
	// Swagger関連のインポート
	_ "github.com/shuto.sawaki/elmo-project/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	// Swagger UIのルートを追加
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ヘルスチェックエンドポイント
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
			"version": "1.0.0",
		})
	})

	// --- ルーティングの設定 ---
	// GETとPOSTのようにHTTPメソッドごとに明確にルートを定義できます。
	// これにより、ハンドラー内のswitch文が不要になります。
	router.GET("/rooms", roomHandler.GetRooms)
	router.POST("/rooms", roomHandler.CreateRoom)

	// URL内の可変部分を :id のようにコロンで指定できます。
	// これを「URLパラメータ」と呼びます。
	router.GET("/rooms/:id", roomHandler.GetRoomByID)
	router.GET("/rooms/:id/start", roomHandler.StartRoom)
	router.PUT("/rooms/:id/status", roomHandler.UpdateRoomStatus)
	router.GET("/rooms/:id/result", roomHandler.GetRoomResult)
	router.POST("/rooms/:id/conclusion", roomHandler.SaveConclusion)
	router.POST("/rooms/:id/sorena", roomHandler.HandleSorena)
	router.GET("/rooms/:id/summary", roomHandler.CreateSummary)
	router.GET("/rooms/:id/open", roomHandler.OpenRoom)

	router.POST("/users", userHandler.CreateUser)

	router.GET("/participants", participantHandler.GetParticipants)
	router.POST("/participants", participantHandler.AddParticipant)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("サーバー起動: http://localhost:%s", port)
	log.Printf("Swagger UI: http://localhost:%s/swagger/index.html", port)
	log.Printf("ヘルスチェック: http://localhost:%s/health", port)
	// ★ Ginのルーターでサーバーを起動
	router.Run(":" + port)
}