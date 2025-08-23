package models

// SorenaCount 「それな」カウント情報
type SorenaCount struct {
	RoomID string `json:"room_id" example:"room123" description:"会議室のID"`
	UserID string `json:"user_id" example:"user123" description:"ユーザーのID"`
	Count  int    `json:"count" example:"3" description:"「それな」の数"`
}

// SorenaRequest 「それな」処理リクエスト
type SorenaRequest struct {
	UserID string `json:"user_id" example:"user123" description:"ユーザーのID"`
	Count  int    `json:"count" example:"1" description:"追加する「それな」の数"`
}
