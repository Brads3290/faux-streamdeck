package streamdeck_server

import (
	"encoding/json"
	"errors"
	"faux-streamdeck/streamdeck_server/config"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"
)

var idExtractor = regexp.MustCompile(`/command/(\w{8}-\w{4}-\w{4}-\w{4}-\w{12})`)
var chanQueue chan string

// StartServer starts the http server to listen for commands from the client
func StartServer() (chan error, error) {
	var err error

	// Pre-server initialization
	RegisterHandlers()
	chanQueue = StartQueueThread()

	// Start the server on a separate thread. chErr will be given an error after the server returns.
	chErr := make(chan error)
	go func() {
		err := http.ListenAndServe(config.General.Server.ListenOn, nil)
		chErr <- err
	}()

	// Wait for a short period of time before continuing, to make sure the server starts OK.
	// If it doesn't, just return the error.
	select {
	case err = <-chErr:
		log.Println("Server failed to start:", err)
		break
	case <-time.After(250 * time.Millisecond):
		log.Println("Server running on", config.General.Server.ListenOn)
		break
	}

	return chErr, err
}

// RegisterHandlers registers the handlers for the HTTP server
func RegisterHandlers() {
	http.HandleFunc("/commands", GetCommands)
	http.HandleFunc("/command/", QueueCommand)
}

// GetCommands handles HTTP requests on the /commands endpoint. It returns a uniform JSON response
// by wrapping it's logic in a call to unifyResponse
func GetCommands(w http.ResponseWriter, r *http.Request) {
	unifyResponse(w, func() (interface{}, error) {
		if r.Method != "GET" {
			return nil, errors.New("invalid HTTP method")
		}

		return config.Buttons, nil
	})
}

// QueueCommand handles HTTP requests on the /command/{guid} endpoint. QueueCommand does not execute the command
// logic itself, but queues the command guid for execution by a worker thread.
// It returns a uniform JSON response by wrapping it's logic in a call to unifyResponse
func QueueCommand(w http.ResponseWriter, r *http.Request) {
	unifyResponse(w, func() (interface{}, error) {
		if r.Method != "POST" {
			return nil, errors.New("invalid HTTP method")
		}

		// Path will be in the form:
		// /command/{guid}

		// Extract the guid
		path := r.URL.Path
		matches := idExtractor.FindStringSubmatch(path)
		if len(matches) == 0 {
			return nil, errors.New("Invalid resource path: " + path)
		}

		guid := matches[1]

		// Check if the guid matches any buttons
		found := false
		for _, v := range config.Buttons.Buttons {
			if v.Id == guid {
				found = true
				break
			}
		}

		// Guid matches no buttons; return an error to the client.
		if !found {
			return nil, errors.New("guid does not match any button")
		}

		// Guid matches a button; queue for processing.
		chanQueue <- guid

		return nil, nil
	})
}

// unifyResponse is a wrapper for HTTP endpoint logic, ensuring that responses are
// always returned in a uniform manner. unifyResponse also catches any panics from
// the endpoint logic and returns those as errors.
func unifyResponse(w http.ResponseWriter, action func() (interface{}, error)) {
	var data interface{}
	var err error

	func() {
		defer func() {
			if r := recover(); r != nil {
				err = errors.New(fmt.Sprintf("%v", r))
			}
		}()

		data, err = action()
	}()

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

// writeResponse marshals response data to JSON and writes it to the response stream
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
