package models

// LogEntry は会話ログ1つ分のデータを表します。
type LogEntry struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

// SummaryRequest は /summary エンドポイントが受け取るリクエストボディの構造を表します。
type SummaryRequest struct {
	Logs []LogEntry `json:"logs"`
}