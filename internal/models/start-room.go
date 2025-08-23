package models

// StartRoomResponse 会議室開始時のレスポンス
type StartRoomResponse struct {
	InitialQuestion string        `json:"initial_question" example:"今日の議題について何か質問はありますか？" description:"AIが生成した初期質問"`
	RoomInfo        RoomInfo      `json:"room_info" description:"会議室の基本情報"`
	Participants    []ParticipantUser `json:"participants" description:"参加者の一覧"`
}

// RoomInfo 会議室の情報
type RoomInfo struct {
	RoomID string `json:"room_id" example:"room123" description:"会議室のID"`
	Title  string `json:"title" example:"週次ミーティング" description:"会議室のタイトル"`
	Status string `json:"status" example:"inprogress" description:"会議室のステータス"`
}
