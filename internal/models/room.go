package models

type Room struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Conclusion  string `json:"conclusion,omitempty"`
}

type ConclusionRequest struct {
	Conclusion string `json:"conclusion"`
}
