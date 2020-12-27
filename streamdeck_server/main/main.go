package main

import (
	"faux-streamdeck/streamdeck_server/config"
	"faux-streamdeck/streamdeck_server/streamdeck_server"
	"flag"
	"log"
)

func main() {
	configDir := flag.String("c", "", "The configuration directory")
	flag.Parse()

	err := config.Load(*configDir)
	if err != nil {
		panic(err)
	}

	// Start the server that handles client requests/commands
	chErr, err := streamdeck_server.StartServer()
	if err != nil {
		return
	}

	// TODO: We may want to do stuff in here

	err = <-chErr
	if err != nil {
		log.Println("Webserver finished with: ", err)
	}

	return
}
