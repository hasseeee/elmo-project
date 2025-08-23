package models

import "time"

// ResultRoomInfo リザルト画面の部屋情報を表します
type ResultRoomInfo struct {
	RoomID string `json:"room_id" example:"room123" description:"会議室のID"`
	Title  string `json:"title" example:"週次ミーティング" description:"会議室のタイトル"`
}

// SorenaParticipant リザルト画面の参加者ごとの「それな」数を表します
type SorenaParticipant struct {
	UserID   string `json:"user_id" example:"user123" description:"ユーザーのID"`
	UserName string `json:"user_name" example:"田中太郎" description:"ユーザーの名前"`
	Count    int    `json:"count" example:"5" description:"「それな」の数"`
}

// SorenaSummary リザルト画面の「それな」集計情報を表します
type SorenaSummary struct {
	TotalCount   int                 `json:"total_count" example:"15" description:"「それな」の総数"`
	Participants []SorenaParticipant `json:"participants" description:"参加者ごとの「それな」集計"`
}

// ChatLog リザルト画面のチャットログ一件を表します
type ChatLog struct {
	LogID     int        `json:"log_id" example:"1" description:"ログの一意のID"`
	UserID    *string    `json:"user_id" example:"user123" description:"ユーザーのID（Null許容）"`
	Message   string     `json:"message" example:"良いアイデアですね" description:"チャットメッセージ"`
	IsSummary bool       `json:"is_summary" example:"false" description:"要約メッセージかどうか"`
	Timestamp time.Time  `json:"timestamp" example:"2024-01-01T10:00:00Z" description:"タイムスタンプ"`
}

// RoomResultResponse リザルト画面APIの完全なレスポンスボディを表します
type RoomResultResponse struct {
	RoomInfo      ResultRoomInfo  `json:"room_info" description:"会議室の基本情報"`
	SorenaSummary SorenaSummary   `json:"sorena_summary" description:"「それな」の集計情報"`
	ChatLogs      []ChatLog       `json:"chat_logs" description:"チャットログの一覧"`
}