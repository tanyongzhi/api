package config

import (
	"os"

	"github.com/HackIllinois/api/common/configloader"
)

var ROOM_PORT string

var ROOM_DB_HOST string
var ROOM_DB_NAME string

func Initialize() error {
	cfg_loader, err := configloader.Load(os.Getenv("HI_CONFIG"))

	ROOM_PORT, err = cfg_loader.Get("ROOM_PORT")

	if err != nil {
		return err
	}

	ROOM_DB_HOST, err = cfg_loader.Get("ROOM_DB_HOST")

	if err != nil {
		return err
	}

	ROOM_DB_NAME, err = cfg_loader.Get("ROOM_DB_NAME")

	if err != nil {
		return err
	}

	return err
}
