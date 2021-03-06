package config

import (
	"encoding/xml"
	"github.com/google/uuid"
	"io/ioutil"
	"os"
	"path/filepath"
)

var Buttons ButtonListSchema
var General StreamdeckConfigSchema

const (
	CONFIG_COMMANDS = "buttons.config"
	CONFIG_GENERAL  = "streamdeck.config"
)

// Load loads the configuration files from the configuration directory that it's given.
func Load(directoryPath string) error {

	// Check the directory path exists
	_, err := os.Stat(directoryPath)
	if err != nil {
		return err
	}

	// Load the command config
	commandsConfigPath := filepath.Join(directoryPath, CONFIG_COMMANDS)
	_, err = os.Stat(commandsConfigPath)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadFile(commandsConfigPath)

	err = xml.Unmarshal(b, &Buttons)
	if err != nil {
		return err
	}

	// Load the faux-streamdeck general config
	streamdeckConfigPath := filepath.Join(directoryPath, CONFIG_GENERAL)
	_, err = os.Stat(streamdeckConfigPath)
	if err != nil {
		return err
	}

	b2, err := ioutil.ReadFile(streamdeckConfigPath)

	err = xml.Unmarshal(b2, &General)
	if err != nil {
		return err
	}

	// Iterate the button list and remove any buttons that don't have a name attribute
	finalButtonsList := make([]Button, 0)
	for _, v := range Buttons.Buttons {
		if v.Name == "" {
			continue
		}

		u, err := uuid.NewRandom()
		if err != nil {
			return err
		}

		v.Id = u.String()
		finalButtonsList = append(finalButtonsList, v)
	}

	Buttons.Buttons = finalButtonsList
	return nil
}
