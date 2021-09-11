package controller

import (
	"encoding/json"
	"net/http"

	"github.com/HackIllinois/api/common/database"
	"github.com/HackIllinois/api/common/errors"
	"github.com/HackIllinois/api/services/room/models"
	"github.com/HackIllinois/api/services/room/service"
	"github.com/gorilla/mux"
)

func SetupController(route *mux.Route) {
	router := route.Subrouter()

	router.HandleFunc("/update/", UpdateRoomOccupancy).Methods("POST")
	router.HandleFunc("/occupancy/{id}/", GetRoomOccupancyById).Methods("GET")
	router.HandleFunc("/occupancy/", GetAllRoomOccupancy).Methods("GET")

}

/*
	Endpoint to update the current occupancy of a room
*/
func UpdateRoomOccupancy(w http.ResponseWriter, r *http.Request) {
	var room_modification models.RoomOccupancyModification
	json.NewDecoder(r.Body).Decode(&room_modification)

	if room_modification.RoomID == "" {
		errors.WriteError(w, r, errors.MalformedRequestError("Must provide room id parameter in request.", "Must provide room id parameter in request."))
		return
	}
}

/*
	Endpoint to get the currency occupancy of a room, by its unique id
*/
func GetRoomOccupancyById(w http.ResponseWriter, r *http.Request) {
	room_id := mux.Vars(r)["id"]

	db_resp, err := service.GetRoomOccupancyById(room_id)

	if err != nil {
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
