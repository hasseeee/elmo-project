package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

var (
	lastNames = []string{
		"佐藤", "鈴木", "高橋", "田中", "伊藤",
		"渡辺", "山本", "中村", "小林", "加藤",
		"吉田", "山田", "佐々木", "山口", "松本",
		"井上", "木村", "林", "斎藤", "清水",
		"山崎", "阿部", "森", "池田", "橋本",
		"山下", "石川", "中島", "前田", "藤田",
		"後藤", "小川", "岡田", "村上", "長谷川",
		"近藤", "石井", "斉藤", "坂本", "遠藤",
		"藤井", "青木", "福田", "三浦", "西村",
		"太田", "松田", "原田", "岡本", "中野",
	}

	firstNames = []string{
		"翔太", "陽菜", "大輝", "美咲", "悠斗",
		"さくら", "悠真", "美羽", "颯太", "葵",
		"蓮", "結衣", "樹", "心愛", "大和",
		"優奈", "陸", "莉子", "悠人", "愛菜",
		"湊", "美優", "朝陽", "結菜", "悠太",
		"杏", "陽太", "優花", "海斗", "心",
		"翼", "彩花", "大翔", "美月", "瑛太",
		"桜", "悠", "綾乃", "大地", "千尋",
		"遥斗", "美桜", "悠希", "七海", "大輔",
		"優衣", "光", "真央", "海", "詩織",
	}

	roomTypes = []string{
		"雑談", "趣味", "仕事", "勉強", "ゲーム",
		"音楽", "映画", "スポーツ", "料理", "旅行",
		"アニメ", "マンガ", "読書", "プログラミング", "デザイン",
		"写真", "動画", "イラスト", "DIY", "ガーデニング",
		"ファッション", "美容", "健康", "ペット", "車",
		"バイク", "自転車", "釣り", "キャンプ", "登山",
		"ランニング", "ヨガ", "ダンス", "カラオケ", "楽器",
		"コレクション", "手芸", "工作", "お菓子作り", "カフェ",
		"お酒", "グルメ", "インテリア", "文具", "ガジェット",
		"占い", "アート", "英会話", "資格", "投資",
	}

	roomDescriptions = []string{
		"みんなで楽しく話しましょう！",
		"初心者から上級者まで大歓迎です",
		"情報交換の場として使ってください",
		"気軽に参加してください",
		"一緒に上達しましょう",
		"経験者の方、アドバイスお願いします",
		"新しい発見を共有しましょう",
		"毎日の成長を記録しよう",
		"仲間と一緒に頑張りましょう",
		"楽しみながら学びましょう",
	}
)

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

	// データベースに接続
	connStr := fmt.Sprintf("user=postgres password=%s dbname=elmo-db sslmode=disable", password)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal("データベースへの接続に失敗しました:", err)
	}
	defer db.Close()

	// 乱数のシードを設定
	rand.Seed(time.Now().UnixNano())

	// 既存のデータを削除
	clearTables(db)

	// ランダムなデータを生成して挿入
	insertRandomData(db, 50) // 各テーブルに50個のデータを挿入
}

func clearTables(db *sql.DB) {
	tables := []string{"sorena_counts", "rooms", "users"}
	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			log.Printf("%sテーブルのクリアに失敗しました: %v", table, err)
		} else {
			log.Printf("%sテーブルをクリアしました", table)
		}
	}
}

func getRandomElement(slice []string) string {
	return slice[rand.Intn(len(slice))]
}

func insertRandomData(db *sql.DB, count int) {
	// ユーザーを作成
	userIDs := make([]string, count)
	for i := 0; i < count; i++ {
		id, err := gonanoid.Generate("0123456789abcdefghijklmnopqrstuvwxyz", 8)
		if err != nil {
			log.Printf("ユーザーID生成エラー: %v", err)
			continue
		}
		userIDs[i] = id

		userName := fmt.Sprintf("%s %s", getRandomElement(lastNames), getRandomElement(firstNames))
		_, err = db.Exec(
			"INSERT INTO users (id, user_name) VALUES ($1, $2)",
			id,
			userName,
		)
		if err != nil {
			log.Printf("ユーザー作成エラー: %v", err)
		}
	}
	log.Printf("%d人のユーザーを作成しました", count)

	// 部屋を作成
	roomIDs := make([]string, count)
	for i := 0; i < count; i++ {
		id, err := gonanoid.Generate("0123456789abcdefghijklmnopqrstuvwxyz", 6)
		if err != nil {
			log.Printf("部屋ID生成エラー: %v", err)
			continue
		}
		roomIDs[i] = id

		roomType := getRandomElement(roomTypes)
		description := getRandomElement(roomDescriptions)
		_, err = db.Exec(
			"INSERT INTO rooms (id, title, description, conclusion) VALUES ($1, $2, $3, $4)",
			id,
			fmt.Sprintf("%sの部屋", roomType),
			description,
			"", // 初期値は空文字列
		)
		if err != nil {
			log.Printf("部屋作成エラー: %v", err)
		}
	}
	log.Printf("%d個の部屋を作成しました", count)

	// それなカウントを作成
	for i := 0; i < count; i++ {
		userID := userIDs[rand.Intn(len(userIDs))]
		roomID := roomIDs[rand.Intn(len(roomIDs))]
		count := rand.Intn(100) + 1 // 1から100のランダムな数

		_, err := db.Exec(
			"INSERT INTO sorena_counts (room_id, user_id, count) VALUES ($1, $2, $3)",
			roomID,
			userID,
			count,
		)
		if err != nil {
			log.Printf("それなカウント作成エラー: %v", err)
		}
	}
	log.Printf("%d個のそれなカウントを作成しました", count)
}
