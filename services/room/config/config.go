package config

import (
	"os"

	"github.com/HackIllinois/api/common/configloader"
)

var ROOM_PORT string

func Initialize() error {
	cfg_loader, err := configloader.Load(os.Getenv("HI_CONFIG"))

	ROOM_PORT, err = cfg_loader.Get("ROOM_PORT")

	return err
}
