package models

// LogEntry 会議のログ一件を表す構造体です
type LogEntry struct {
	Content   string `json:"content" example:"プロジェクトの進捗について話し合いました" description:"ログの内容"`
	UserID    string `json:"user_id,omitempty" description:"ログを作成したユーザーのID"`
	IsSummary bool   `json:"is_summary" description:"要約かどうか"`
	CreatedAt string `json:"created_at" description:"ログの作成日時"`
}