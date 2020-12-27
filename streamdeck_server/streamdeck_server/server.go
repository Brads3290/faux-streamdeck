package streamdeck_server

import (
	"encoding/json"
	"errors"
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
	unifyResponse(w, func() (interface{}, error) {
		if r.Method != "GET" {
			return nil, errors.New("invalid HTTP method")
		}

		return config.Buttons, nil
	})
}

func ExecuteCommand(w http.ResponseWriter, r *http.Request) {
	unifyResponse(w, func() (interface{}, error) {
		if r.Method != "POST" {
			return nil, errors.New("invalid HTTP method")
		}

		// Path will be in the form:
		// /command/{guid}
		path := r.URL.Path
		matches := idExtractor.FindStringSubmatch(path)
		if len(matches) == 0 {
			return nil, errors.New("Invalid resource path: " + path)
		}

		guid := matches[1]

		found := false
		for _, v := range config.Buttons.Buttons {
			if v.Id == guid {
				found = true
				break
			}
		}

		if !found {
			return nil, errors.New("guid does not match any button")
		}

		chanQueue <- guid

		return nil, nil
	})
}

func unifyResponse(w http.ResponseWriter, action func() (interface{}, error)) {
	data, err := action()

	jsonResponse := struct {
		Error  bool        `json:"error"`
		Reason string      `json:"reason"`
		Data   interface{} `json:"data"`
	}{
		Error:  false,
		Reason: "",
		Data:   nil,
	}

	if err != nil {
		jsonResponse.Error = true
		jsonResponse.Reason = err.Error()
	} else {
		jsonResponse.Error = false
		jsonResponse.Reason = ""
		jsonResponse.Data = data
	}

	writeResponse(w, jsonResponse)
	return
}

func writeResponse(w http.ResponseWriter, responseData interface{}) {
	b, err := json.Marshal(responseData)
	if err != nil {
		log.Println("ERROR - Failed to marshal json response. Error = ", err)
	}

	_, err = w.Write(b)
	if err != nil {
		log.Println("ERROR - Failed to write to response writer. Error = ", err)
	}

	return
}
