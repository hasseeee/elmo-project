package models

// User ユーザー情報を表す構造体
type User struct {
	ID       string `json:"id" example:"user123" description:"ユーザーの一意のID"`
	UserName string `json:"user_name" example:"田中太郎" description:"ユーザーの名前"`
}
