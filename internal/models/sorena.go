package models

type SorenaCount struct {
	RoomID string `json:"room_id"`
	UserID string `json:"user_id"`
	Count  int    `json:"count"`
}

type SorenaRequest struct {
	UserID string `json:"user_id"`
	Count  int    `json:"count"`
}
