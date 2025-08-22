package models

type Participant struct {
	RoomID string `json:"room_id"`
	UserID string `json:"user_id"`
}

type ParticipantRequest struct {
	RoomID string `json:"room_id"`
	UserID string `json:"user_id"`
}

type ParticipantUser struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
}

type ParticipantsResponse struct {
	RoomID string           `json:"room_id"`
	Users  []ParticipantUser `json:"users"`
}
