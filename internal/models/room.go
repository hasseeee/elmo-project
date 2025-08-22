package models

type Room struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Conclusion  string `json:"conclusion,omitempty"`
	Status      string `json:"status,omitempty"`
	InitialQuestion  string `json:"initial_question,omitempty"` // 問いかけを追加

}

type ConclusionRequest struct {
	Conclusion string `json:"conclusion"`
}