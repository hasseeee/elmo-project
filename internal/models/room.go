package models

// Room 会議室の情報を表す構造体
type Room struct {
	ID          string `json:"id" example:"abc123" description:"会議室の一意のID"`
	Title       string `json:"title" example:"週次ミーティング" description:"会議室のタイトル"`
	Description string `json:"description" example:"今週の進捗確認と来週の計画" description:"会議室の説明"`
	Conclusion  string `json:"conclusion,omitempty" example:"来週までにプロトタイプを完成させる" description:"会議の結論（オプション）"`
	Status      string `json:"status,omitempty" example:"inprogress" description:"会議室のステータス（オプション）"`
	InitialQuestion  string `json:"initial_question,omitempty" example:"今日の議題について何か質問はありますか？" description:"AIが生成した初期質問（オプション）"`
}

// UpdateRoomStatusRequest 会議室のステータス更新リクエスト
type UpdateRoomStatusRequest struct {
	Status string `json:"status" example:"done" description:"更新するステータス"`
}


