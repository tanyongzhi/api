package service

import (
	"fmt"

	"github.com/HackIllinois/api/common/database"
	"github.com/HackIllinois/api/services/room/config"
	"github.com/HackIllinois/api/services/room/models"
)

var db database.Database

/*
	Custom errors
*/
type ErrNegativeRemainingSpaces struct {
	RemainingSpaces int
	RoomID          string
}

func (e *ErrNegativeRemainingSpaces) Error() string {
	return fmt.Sprintf("Invalid operation: room %v only has %v remaining spaces", e.RoomID, e.RemainingSpaces)
}

type ErrExceedRemainingSpaces struct {
	MaxCapacity int
	RoomID      string
}

func (e *ErrExceedRemainingSpaces) Error() string {
	return fmt.Sprintf("Invalid operation:  room %v has a max capacity of %v", e.RoomID, e.MaxCapacity)
}

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
func GetRoomOccupancyById(room_id string) (models.RoomOccupancy, error) {
	query := database.QuerySelector{
		"roomId": room_id,
	}

	var occupancy models.RoomOccupancy
	err := db.FindOne(OCCUPANCY_COLLECTION, query, &occupancy)

	return occupancy, err
}

/*
	Fetches occupancy values corresponding to all roomIDs
*/
func GetAllRoomOccupancy() ([]models.RoomOccupancy, error) {
	query := database.QuerySelector{}

	var occupancy []models.RoomOccupancy
	err := db.FindAll(OCCUPANCY_COLLECTION, query, &occupancy)

	return occupancy, err
}

/*
	Writes the new occupancy value corresponding to the respective roomId to the database.

	Also performs validity checks if the new remaining spaces value is invalid
*/
func UpdateRoomOccupancy(room_id string, change_in_remaining_spaces int) error {
	// read from db
	db_resp, err := GetRoomOccupancyById(room_id)
	if err != nil {
		return err
	}

	// validity checks
	new_remaining_spaces := db_resp.RemainingSpaces - change_in_remaining_spaces
	err = checkRoomSpaceValid(db_resp, new_remaining_spaces)
	if err != nil {
		return err
	}

	// checks are ok, write to db now
	selector := database.QuerySelector{
		"roomId": room_id,
	}
	err = db.Update(OCCUPANCY_COLLECTION, selector, &models.RoomOccupancy{
		RoomID:          room_id,
		RemainingSpaces: new_remaining_spaces,
		MaxCapacity:     db_resp.MaxCapacity,
	})
	return err
}

/*
	Validity checks on the remaining spaces in a room

	Returns (err msg, validity check success/fail result)
*/
func checkRoomSpaceValid(db_resp models.RoomOccupancy, new_remaining_spaces int) error {
	if new_remaining_spaces < 0 {
		return &ErrNegativeRemainingSpaces{
			RemainingSpaces: db_resp.RemainingSpaces,
			RoomID:          db_resp.RoomID,
		}
	}
	if new_remaining_spaces > db_resp.MaxCapacity {
		return &ErrExceedRemainingSpaces{
			MaxCapacity: db_resp.MaxCapacity,
			RoomID:      db_resp.RoomID,
		}
	}

	return nil
}
