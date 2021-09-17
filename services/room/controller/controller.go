package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/HackIllinois/api/common/database"
	"github.com/HackIllinois/api/common/errors"
	"github.com/HackIllinois/api/services/room/models"
	"github.com/HackIllinois/api/services/room/service"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	remaining_spaces_metrics = *promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "room_remaining_spaces",
		Help: "Number of remaining spaces for each room",
	}, []string{"roomID"})
)

func SetupController(route *mux.Route) {
	emitOccupancyCounts()

	router := route.Subrouter()
	router.HandleFunc("/update/", UpdateRoomOccupancy).Methods("POST")
	router.HandleFunc("/occupancy/{id}/", GetRoomOccupancyById).Methods("GET")
	router.HandleFunc("/occupancy/", GetAllRoomOccupancy).Methods("GET")
	router.Handle("/graph", promhttp.Handler()).Methods("GET")

}

func emitOccupancyCounts() {
	QUERY_TIME := 2 * time.Second
	go func() {
		for {
			db_resp, err := service.GetAllRoomOccupancy()
			if err == nil {
				for _, room := range db_resp {
					remaining_spaces_metrics.WithLabelValues(room.RoomID).Set(float64(room.RemainingSpaces))
				}
			}

			time.Sleep(QUERY_TIME)
		}
	}()
}

/*
	Endpoint to update the current occupancy of a room
*/
func UpdateRoomOccupancy(w http.ResponseWriter, r *http.Request) {
	var room_modification models.RoomOccupancyModification
	json.NewDecoder(r.Body).Decode(&room_modification)

	var room_id string
	if room_id = room_modification.RoomID; room_id == "" {
		errors.WriteError(w, r, errors.MalformedRequestError("Must provide room id parameter in request.", "Must provide room id parameter in request."))
		return
	}

	// retrieve pertinent record from db
	db_resp, err := service.GetRoomOccupancyById(room_id)
	if err != nil {
		message, http_err_type, http_status_code := generateDbHttpErr(err)
		errors.WriteError(w, r, errors.ApiError{Status: http_status_code, Type: http_err_type, Message: message, RawError: err.Error()})
		return
	}

	// validity checks
	new_remaining_spaces := db_resp.RemainingSpaces - room_modification.NumPeople
	error_msg, is_valid := checkRoomSpaceValid(db_resp, new_remaining_spaces)
	if !is_valid {
		errors.WriteError(w, r, errors.ApiError{Status: http.StatusForbidden, Type: "INVALID_OPERATION", Message: error_msg, RawError: error_msg})
		return
	}

	// checks are ok, write to db now
	err = service.UpdateRoomOccupancy(room_id, new_remaining_spaces, db_resp.MaxCapacity)
	if err != nil {
		message, http_err_type, http_status_code := generateDbHttpErr(err)
		errors.WriteError(w, r, errors.ApiError{Status: http_status_code, Type: http_err_type, Message: message, RawError: err.Error()})
		return
	}

	// final read, then return
	db_resp, err = service.GetRoomOccupancyById(room_id)
	if err != nil {
		message, http_err_type, http_status_code := generateDbHttpErr(err)
		errors.WriteError(w, r, errors.ApiError{Status: http_status_code, Type: http_err_type, Message: message, RawError: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(db_resp)
}

/*
	Validity checks on the remaining spaces in a room

	Returns (err msg, validity check success/fail result)
*/
func checkRoomSpaceValid(db_resp models.RoomOccupancy, new_remaining_spaces int) (string, bool) {
	var msg string = ""

	if new_remaining_spaces < 0 {
		msg = fmt.Sprintf("Invalid operation: only %v remaining spaces left", db_resp.RemainingSpaces)
		return msg, false
	}
	if new_remaining_spaces > db_resp.MaxCapacity {
		msg = fmt.Sprintf("Invalid operation: the max capacity of room %v is %v", db_resp.RoomID, db_resp.MaxCapacity)
		return msg, false
	}

	return msg, true
}

/*
	Endpoint to get the currency occupancy of a room, by its unique id
*/
func GetRoomOccupancyById(w http.ResponseWriter, r *http.Request) {
	room_id := mux.Vars(r)["id"]

	db_resp, err := service.GetRoomOccupancyById(room_id)
	if err != nil {
		message, http_err_type, http_status_code := generateDbHttpErr(err)
		errors.WriteError(w, r, errors.ApiError{Status: http_status_code, Type: http_err_type, Message: message, RawError: err.Error()})
	}

	json.NewEncoder(w).Encode(db_resp)
}

/*
	Endpoint to get the currency occupancy of a room
*/
func GetAllRoomOccupancy(w http.ResponseWriter, r *http.Request) {
	db_resp, err := service.GetAllRoomOccupancy()

	if err != nil {
		var message, http_err_type string
		var http_status_code int

		switch err {
		case database.ErrConnection:
			message, http_err_type = "Connection error to database", "CONN_ERR"
			http_status_code = http.StatusInternalServerError
		default:
			message, http_err_type = "Unknown error", "UNKNOWN_ERROR"
			http_status_code = http.StatusInternalServerError
		}

		errors.WriteError(w, r, errors.ApiError{Status: http_status_code, Type: http_err_type, Message: message, RawError: err.Error()})
	}

	json.NewEncoder(w).Encode(db_resp)
}

/*
	Helper function to form HTTP error message, code, and error type
*/
func generateDbHttpErr(err error) (string, string, int) {
	var message, http_err_type string
	var http_status_code int

	switch err {
	case database.ErrNotFound:
		message, http_err_type = "Room ID does not exist", "NOT_FOUND"
		http_status_code = http.StatusNotFound
	case database.ErrConnection:
		message, http_err_type = "Connection error to database", "CONN_ERR"
		http_status_code = http.StatusInternalServerError
	default:
		message, http_err_type = "Unknown error", "UNKNOWN_ERROR"
		http_status_code = http.StatusInternalServerError
	}

	return message, http_err_type, http_status_code
}
