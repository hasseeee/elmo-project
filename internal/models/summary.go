package models

// SummaryRequest /summary エンドポイントが受け取るリクエストボディの構造を表します
type SummaryRequest struct {
	Logs []LogEntry `json:"logs" description:"要約対象のログエントリの一覧"`
}