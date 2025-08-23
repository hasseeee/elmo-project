package models

// Participant 参加者情報を表す構造体
type Participant struct {
	RoomID string `json:"room_id" example:"room123" description:"会議室のID"`
	UserID string `json:"user_id" example:"user123" description:"ユーザーのID"`
}

// ParticipantRequest 参加者追加リクエスト
type ParticipantRequest struct {
	RoomID string `json:"room_id" example:"room123" description:"会議室のID"`
	UserID string `json:"user_id" example:"user123" description:"ユーザーのID"`
}

// ParticipantUser 参加者ユーザー情報
type ParticipantUser struct {
	ID           string `json:"id" example:"user123" description:"ユーザーの一意のID"`
	Name         string `json:"name" example:"田中太郎" description:"ユーザーの名前"`
}

// ParticipantsResponse 参加者一覧レスポンス
type ParticipantsResponse struct {
	RoomID string           `json:"room_id" example:"room123" description:"会議室のID"`
	Users  []ParticipantUser `json:"users" description:"参加者ユーザーの一覧"`
}
