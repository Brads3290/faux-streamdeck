package streamdeck_server

import (
	"net/http"
	"streamdeck_server/config"
)

// StartServer starts the http server to listen to the
func StartServer() chan error {

	http.HandleFunc("/commands", GetCommands)
	http.HandleFunc("/command", ExecuteCommand)

	chErr := make(chan error)
	go func() {
		err := http.ListenAndServe(config.General.Server.ListenOn, nil)
		chErr <- err
	}()

	return chErr
}

func GetCommands(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}

	return
}

func ExecuteCommand(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	return
}





