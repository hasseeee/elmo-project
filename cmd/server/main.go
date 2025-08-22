package main

import (
	"context"
	"log"
	"net/http"

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

	// ハンドラーの初期化
	roomHandler := handlers.NewRoomHandler(database, aiGenerator)
	userHandler := handlers.NewUserHandler(database)
	participantHandler := handlers.NewParticipantHandler(database)

	// --- ルーティングの設定 ---
	mux := http.NewServeMux()

	// Room Handlers
	mux.HandleFunc("/rooms", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			roomHandler.GetRooms(w, r)
		case http.MethodPost:
			roomHandler.CreateRoom(w, r)
		default:
			http.Error(w, "サポートされていないメソッドです", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/rooms/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.HasSuffix(path, "/start") {
			roomHandler.StartRoom(w, r)
		} else if strings.HasSuffix(path, "/conclusion") {
			roomHandler.SaveConclusion(w, r)
		} else if strings.HasSuffix(path, "/sorena") {
			roomHandler.HandleSorena(w, r)
		} else {
			roomHandler.GetRoomByID(w, r)
		}
	})

	// User Handlers
	mux.HandleFunc("/users", userHandler.CreateUser)

	// Participant Handlers
	mux.HandleFunc("/participants", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			participantHandler.GetParticipants(w, r)
		case http.MethodPost:
			participantHandler.AddParticipant(w, r)
		default:
			http.Error(w, "サポートされていないメソッドです", http.StatusMethodNotAllowed)
		}
	})


	log.Println("サーバー起動: http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("サーバーの起動に失敗しました:", err)
	}
}