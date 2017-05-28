package main

import (
	"bitbucket.org/rawfish-dev/wedding-rsvp/server/api"
	"bitbucket.org/rawfish-dev/wedding-rsvp/server/config"
)

func main() {
	loadedConfig := config.Load()

	reactReduxBasicsAPI := api.NewAPI(*loadedConfig)

	reactReduxBasicsAPI.Run()
}
