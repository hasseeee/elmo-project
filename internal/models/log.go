package models

import "time"

type LogEntry struct {
	ID        int       `json:"id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	IsSummary bool      `json:"is_summary"`
	CreatedAt time.Time `json:"created_at"`
}