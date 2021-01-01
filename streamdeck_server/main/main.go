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

	// Load the configuration. This is loaded from a directory because there will be multiple
	// configuration files
	err := config.Load(*configDir)
	if err != nil {
		panic(err)
	}

	// Start the server that handles client requests/commands
	chErr, err := streamdeck_server.StartServer()
	if err != nil {
		log.Fatal("Could not start server: ", err)
	}

	// TODO: Listen for interrupt signals and gracefully shutdown the server

	err = <-chErr
	if err != nil {
		log.Println("Webserver finished with: ", err)
	}

	return
}
