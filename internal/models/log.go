package models

// LogEntry 会議のログ一件を表す構造体です
type LogEntry struct {
	Content string `json:"content" example:"プロジェクトの進捗について話し合いました" description:"ログの内容"`
	// 他にもタイムスタンプなどの情報が必要であれば、ここに追加します
