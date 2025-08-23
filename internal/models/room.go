package models

type Room struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Conclusion  string `json:"conclusion,omitempty"`
	Status      string `json:"status,omitempty"`
}

type UpdateRoomStatusRequest struct {
	Status string `json:"status"`
}
