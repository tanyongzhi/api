package models

type RoomOccupancy struct {
	RoomID          string `bson:"roomId" json:"roomId"`
	RemainingSpaces int    `bson:"remainingSpaces" json:"remainingSpaces"`
	MaxCapacity     int    `bson:"maxCapacity" json:"maxCapacity"`
}
