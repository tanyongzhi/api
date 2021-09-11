package models

type RoomOccupancyModification struct {
	RoomID    string `json:"roomId"`
	NumPeople int    `json:"numPeople"`
}
