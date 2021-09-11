package services

import (
	"net/http"

	"github.com/HackIllinois/api/gateway/config"
	"github.com/arbor-dev/arbor"
)

var RoomRoutes = arbor.RouteCollection{
	arbor.Route{
		"UpdateRoomOccupancy",
		"POST",
		"/room/update/",
		UpdateRoomOccupancy,
	},
	arbor.Route{
		"GetRoomOccupancyById",
		"GET",
		"/room/occupancy/{id}/",
		GetRoomOccupancyById,
	},
	arbor.Route{
		"GetAllRoomOccupancy",
		"GET",
		"/room/occupancy/",
		GetAllRoomOccupancy,
	},
}

func UpdateRoomOccupancy(w http.ResponseWriter, r *http.Request) {
	arbor.POST(w, config.ROOM_SERVICE+r.URL.String(), InfoFormat, "", r)
}

func GetRoomOccupancyById(w http.ResponseWriter, r *http.Request) {
	arbor.GET(w, config.ROOM_SERVICE+r.URL.String(), InfoFormat, "", r)
}

func GetAllRoomOccupancy(w http.ResponseWriter, r *http.Request) {
	arbor.GET(w, config.ROOM_SERVICE+r.URL.String(), InfoFormat, "", r)
}
