package models

// SummaryRequest は /summary エンドポイントが受け取るリクエストボディの構造を表します。
type SummaryRequest struct {
	Logs []LogEntry `json:"logs"`
}