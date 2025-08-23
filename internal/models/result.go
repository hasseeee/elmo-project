package models

import "time"

// ResultRoomInfo はリザルト画面の部屋情報を表します。
type ResultRoomInfo struct {
	RoomID string `json:"room_id"`
	Title  string `json:"title"`
}

// SorenaParticipant はリザルト画面の参加者ごとの「それな」数を表します。
type SorenaParticipant struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
	Count    int    `json:"count"`
}

// SorenaSummary はリザルト画面の「それな」集計情報を表します。
type SorenaSummary struct {
	TotalCount   int                 `json:"total_count"`
	Participants []SorenaParticipant `json:"participants"`
}

// ChatLog はリザルト画面のチャットログ一件を表します。
type ChatLog struct {
	LogID     int        `json:"log_id"`
	UserID    *string    `json:"user_id"` // Null許容のためポインタ型
	Message   string     `json:"message"`
	IsSummary bool       `json:"is_summary"`
	Timestamp time.Time  `json:"timestamp"`
}

// RoomResultResponse はリザルト画面APIの完全なレスポンスボディを表します。
type RoomResultResponse struct {
	RoomInfo      ResultRoomInfo  `json:"room_info"`
	SorenaSummary SorenaSummary   `json:"sorena_summary"`
	ChatLogs      []ChatLog       `json:"chat_logs"`
}