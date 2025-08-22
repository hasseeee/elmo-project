package models

// StartRoomResponse represents the response for starting a room
type StartRoomResponse struct {
	InitialQuestion string        `json:"initial_question"`
	RoomInfo        RoomInfo      `json:"room_info"`
	Participants    []ParticipantUser `json:"participants"`
}

// RoomInfo represents information about a room
type RoomInfo struct {
	RoomID string `json:"room_id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}
