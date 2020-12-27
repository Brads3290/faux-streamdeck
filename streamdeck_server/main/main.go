package main

import (
	"flag"
	"log"
	"streamdeck_server/config"
	"streamdeck_server/streamdeck_server"
)

func main() {
	configDir := flag.String("c", "", "The configuration directory")
	flag.Parse()

	err := config.Load(*configDir)
	if err != nil {
		panic(err)
	}

	// Start the server that handles client requests/commands
	chErr := streamdeck_server.StartServer()

	// TODO: We may want to do stuff in here

	err = <- chErr
	if err != nil {
		log.Println("Webserver finished with: ", err)
	}

	return
}
