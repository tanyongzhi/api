package controller

import (
	"net/http"

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
}

/*
	Endpoint to get the currency occupancy of a room, by its unique id
*/
func GetRoomOccupancyById(w http.ResponseWriter, r *http.Request) {
}

/*
	Endpoint to get the currency occupancy of a room
*/
func GetAllRoomOccupancy(w http.ResponseWriter, r *http.Request) {
}
