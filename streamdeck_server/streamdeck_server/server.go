package streamdeck_server

import (
	"encoding/json"
	"faux-streamdeck/streamdeck_server/config"
	"log"
	"net/http"
	"time"
)

// StartServer starts the http server to listen to the
func StartServer() (chan error, error) {

	http.HandleFunc("/commands", GetCommands)
	http.HandleFunc("/command", ExecuteCommand)

	chErr := make(chan error)
	go func() {
		err := http.ListenAndServe(config.General.Server.ListenOn, nil)
		chErr <- err
	}()

	var err error = nil
	select {
	case err = <-chErr:
		log.Println("Server failed to start:", err)
		break
	case <-time.After(1 * time.Second):
		log.Println("Server running on", config.General.Server.ListenOn)
		break
	}

	return chErr, err
}

func GetCommands(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}

	jsonBytes, err := json.Marshal(config.Buttons)
	if err != nil {
		_, err = w.Write([]byte("error"))

		if err != nil {
			log.Println("ERROR: ", err)
		}
	}

	_, err = w.Write(jsonBytes)
	if err != nil {
		log.Println("ERROR: ", err)
	}

	return
}

func ExecuteCommand(_ http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	return
}
