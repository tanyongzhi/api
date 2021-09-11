package service

import (
	"github.com/HackIllinois/api/common/database"
	"github.com/HackIllinois/api/services/room/config"
	"github.com/HackIllinois/api/services/room/models"
)

var db database.Database

const OCCUPANCY_COLLECTION string = "occupancy"

/*
	Initialize DB connections
*/
func Initialize() error {
	if db != nil {
		db.Close()
		db = nil
	}

	var err error
	db, err = database.InitDatabase(config.ROOM_DB_HOST, config.ROOM_DB_NAME)

	if err != nil {
		return err
	}

	return nil
}

/*
	Fetches the occupancy value corresponding to the roomId
*/
func GetRoomOccupancyById(roomId string) (models.RoomOccupancy, error) {
	query := database.QuerySelector{
		"roomId": roomId,
	}

	var occupancy models.RoomOccupancy
	err := db.FindOne(OCCUPANCY_COLLECTION, query, &occupancy)

	return occupancy, err
}

/*
	Fetches occupancy values corresponding to all roomIDs
*/
func GetAllRoomOccupancy() (models.RoomOccupancy, error) {
	query := database.QuerySelector{}

	var occupancy models.RoomOccupancy
	err := db.FindOne(OCCUPANCY_COLLECTION, query, &occupancy)

	return occupancy, err
}

/*
	Writes the new occupancy value corresponding to the respective roomId to the database.

	NOTE: This function does NOT perform any checks regarding the validity of the new occupancy value. Ensure
	that the caller of this function performs all the necessary checks (eg. if negative capacity has been reached)
	before calling this function.
*/
func UpdateRoomOccupancy(roomId int, newOccupancyVal int) error {
	return nil
}
