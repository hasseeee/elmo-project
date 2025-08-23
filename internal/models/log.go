package models

// LogEntry は、会議のログ一件を表す構造体です。
type LogEntry struct {
	Content string
	// 他にもタイムスタンプなどの情報が必要であれば、ここに追加します。
}