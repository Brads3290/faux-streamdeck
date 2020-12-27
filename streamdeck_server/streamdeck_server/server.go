package streamdeck_server

import (
	"encoding/json"
	"faux-streamdeck/streamdeck_server/config"
	"log"
	"net/http"
	"regexp"
	"time"
)

var idExtractor *regexp.Regexp
var chanQueue chan string

// StartServer starts the http server to listen to the
func StartServer() (chan error, error) {
	var err error

	// Pre-server initialization
	CreateRegex()
	RegisterHandlers()
	chanQueue = StartQueueThread()

	// Start the server on a separate thread
	chErr := make(chan error)
	go func() {
		err := http.ListenAndServe(config.General.Server.ListenOn, nil)
		chErr <- err
	}()

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

func RegisterHandlers() {
	http.HandleFunc("/commands", GetCommands)
	http.HandleFunc("/command/", ExecuteCommand)
}

func CreateRegex() {
	var err error

	idExtractor, err = regexp.Compile(`/command/(\w{8}-\w{4}-\w{4}-\w{4}-\w{12})`)
	if err != nil {
		panic(err)
	}
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

	// Path will be in the form:
	// /command/{guid}
	path := r.URL.Path
	matches := idExtractor.FindStringSubmatch(path)
	if len(matches) == 0 {
		log.Println("Invalid request:", path)
		return
	}

	guid := matches[1]
	chanQueue <- guid

	return
}
