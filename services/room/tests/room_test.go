package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/HackIllinois/api/common/database"
	"github.com/HackIllinois/api/services/room/config"
	"github.com/HackIllinois/api/services/room/models"
	"github.com/HackIllinois/api/services/room/service"
)

var db database.Database

const (
	TEST_ROOM_ID_1 = "testroomid1"
	TEST_ROOM_ID_2 = "testroomid2"
	MAX_CAP_1      = 10
	MAX_CAP_2      = 20
)

func TestMain(m *testing.M) {
	err := config.Initialize()

	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)

	}

	err = service.Initialize()

	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}

	db, err = database.InitDatabase(config.ROOM_DB_HOST, config.ROOM_DB_NAME)

	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}

	return_code := m.Run()

	os.Exit(return_code)
}

/*
	Initialize database with test room info
*/
func SetupTestDB(t *testing.T) {
	err := db.Insert("occupancy", &models.RoomOccupancy{
		RoomID:          TEST_ROOM_ID_1,
		RemainingSpaces: MAX_CAP_1,
		MaxCapacity:     MAX_CAP_1,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = db.Insert("occupancy", &models.RoomOccupancy{
		RoomID:          TEST_ROOM_ID_2,
		RemainingSpaces: MAX_CAP_2 - 5,
		MaxCapacity:     MAX_CAP_2,
	})
	if err != nil {
		t.Fatal(err)
	}
}

/*
	Drop test db
*/
func CleanupTestDB(t *testing.T) {
	err := db.DropDatabase()

	if err != nil {
		t.Fatal(err)
	}
}

/*
	Test successful updates with people entering a room
*/
func TestRoomSpacesUpdateEnterSuccess(t *testing.T) {
	SetupTestDB(t)

	test_update_num := 5
	err := service.UpdateRoomOccupancy(TEST_ROOM_ID_1, test_update_num)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := service.GetRoomOccupancyById(TEST_ROOM_ID_1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.RemainingSpaces != MAX_CAP_1-test_update_num {
		t.Errorf("Wrong remaining spaces. Expected %v, got %v", MAX_CAP_1-test_update_num, resp.RemainingSpaces)
	}

	CleanupTestDB(t)
}

/*
	Test successful updates with people leaving a room
*/
func TestRoomSpacesUpdateLeavingSuccess(t *testing.T) {
	SetupTestDB(t)

	test_update_num := -5
	err := service.UpdateRoomOccupancy(TEST_ROOM_ID_2, test_update_num)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := service.GetRoomOccupancyById(TEST_ROOM_ID_2)
	if err != nil {
		t.Fatal(err)
	}
	if resp.RemainingSpaces != MAX_CAP_2 {
		t.Errorf("Wrong remaining spaces. Expected %v, got %v", MAX_CAP_2, resp.RemainingSpaces)
	}

	CleanupTestDB(t)
}

/*
	Test update where remaining seats > max capacity
*/
func TestOverflowUpdate(t *testing.T) {
	SetupTestDB(t)

	test_update_num := -5 - 1
	err := service.UpdateRoomOccupancy(TEST_ROOM_ID_2, test_update_num)
	if err == nil {
		t.Errorf("Expected error in operation overflow")
	}
	if _, ok := err.(*service.ErrExceedRemainingSpaces); !ok {
		t.Errorf("Expected error of type service.ErrExceedRemainingSpaces, but got type %T", err)
	}

	CleanupTestDB(t)
}

/*
	Test update where remaining seats < 0
*/
func TestUnderflowUpdate(t *testing.T) {
	SetupTestDB(t)

	test_update_num := MAX_CAP_1 + 1
	err := service.UpdateRoomOccupancy(TEST_ROOM_ID_1, test_update_num)
	if err == nil {
		t.Errorf("Expected error in operation overflow")
	}
	if _, ok := err.(*service.ErrNegativeRemainingSpaces); !ok {
		t.Errorf("Expected error of type service.ErrNegativeRemainingSpaces, but got type %T", err)
	}

	CleanupTestDB(t)
}

/*
	Test get occupancy for all rooms
*/
func TestGetAllRoomOccupancy(t *testing.T) {
	SetupTestDB(t)

	rooms, err := service.GetAllRoomOccupancy()
	if err != nil {
		t.Fatal(err)
	}

	room_id_set := make(map[string]bool)
	for _, room := range rooms {
		room_id_set[room.RoomID] = true
	}

	if len(room_id_set) != 2 {
		t.Errorf("Expected length %v, got length %v instead", 2, len(room_id_set))
	}
	if in_set, ok := room_id_set[TEST_ROOM_ID_1]; !ok || !in_set {
		t.Errorf("%v is not returned", TEST_ROOM_ID_1)
	}
	if in_set, ok := room_id_set[TEST_ROOM_ID_2]; !ok || !in_set {
		t.Errorf("%v is not returned", TEST_ROOM_ID_2)
	}

	CleanupTestDB(t)
}
