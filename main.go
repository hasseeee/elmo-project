package main

import (
	"encoding/json" // GoのデータをJSON形式に変換するための道具
	"log"           // エラーなどを画面に表示するための道具
	"net/http"      // Webサーバーを動かすための道具
)

// Room 掲示板の部屋の情報をまとめるための「設計図」
type Room struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// アプリが持っている掲示板の部屋一覧（今回は練習なので、プログラム内に直接書いておきます）
var rooms = []Room{
	{ID: 1, Title: "仕事の悩み相談室", Description: "職場の人間関係やキャリアについて話しましょう"},
	{ID: 2, Title: "恋愛の悩み相談室", Description: "恋愛に関する悩みを気軽に相談してください"},
	{ID: 3, Title: "勉強の悩み相談室", Description: "勉強法や進路について情報交換しましょう"},
}

// 「/rooms」にアクセスが来た時に動く関数
func getRoomsHandler(w http.ResponseWriter, r *http.Request) {
	// 1. ヘッダーを設定する
	// これから送るデータは「JSON形式」ですよ、と相手に伝えるためのおまじない
	w.Header().Set("Content-Type", "application/json")

	// 2. Goのデータ（rooms）をJSON形式のデータに変換する
	// json.NewEncoder(w)で、レスポンスに書き込むための準備
	// .Encode(rooms)で、実際にrooms変数をJSONに変換して書き込む
	err := json.NewEncoder(w).Encode(rooms)
	if err != nil {
		// もしJSONへの変換でエラーが起きたら、サーバー側の画面にエラー内容を表示する
		log.Println("JSONへの変換に失敗しました:", err)
	}
}

func main() {
	// 1. Webサーバーに「/rooms」という住所（エンドポイント）が指定されたら、
	// 「getRoomsHandler」という関数を動かしてくださいね、とお願いする
	http.HandleFunc("/rooms", getRoomsHandler)

	// 2. サーバーを起動する
	// "localhost:8080"でサーバーを動かします、と宣言
	log.Println("サーバー起動: http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		// もしサーバー起動でエラーが起きたら、強制終了してエラー内容を表示する
		log.Fatal("サーバーの起動に失敗しました:", err)
	}
}