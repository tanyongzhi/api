package room

import (
	"log"

	"github.com/HackIllinois/api/common/apiserver"
	"github.com/HackIllinois/api/services/room/config"
	"github.com/HackIllinois/api/services/room/controller"
	"github.com/gorilla/mux"
)

func Initialize() error {
	err := config.Initialize()

	if err != nil {
		return err

	}

	// err = service.Initialize()

	if err != nil {
		return err
	}

	return nil
}

func Entry() {
	err := Initialize()

	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	controller.SetupController(router.PathPrefix("/room"))

	log.Fatal(apiserver.StartServer(config.ROOM_PORT, router, "room", Initialize))
}
