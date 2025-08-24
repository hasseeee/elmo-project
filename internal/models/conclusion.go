package models

// ConclusionRequest 会議の結論保存リクエスト
type ConclusionRequest struct {
	Conclusion string `json:"conclusion" example:"来週までにプロトタイプを完成させる" description:"保存する結論"`
}