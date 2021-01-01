package streamdeck_server

import (
	"errors"
	"faux-streamdeck/streamdeck_server/config"
	"fmt"
	"github.com/micmonay/keybd_event"
	"log"
	"runtime"
	"strings"
	"time"
)

//goland:noinspection GoSnakeCaseUsage
const (
	MAX_COMMAND_QUEUE_SIZE = 8
)

func StartQueueThread() chan string {
	chanQueue := make(chan string, MAX_COMMAND_QUEUE_SIZE)

	go runQueueThread(chanQueue)
	return chanQueue
}

func runQueueThread(chanQueue chan string) {
	for {
		guid, ok := <-chanQueue
		if !ok {
			break
		}

		go runCommandOnThread(guid)
	}
}

func runCommandOnThread(commandGuid string) {
	var button *config.Button = nil

	for i, v := range config.Buttons.Buttons {
		if v.Id == commandGuid {
			button = &config.Buttons.Buttons[i]
			break
		}
	}

	if button == nil {
		log.Println("WARN - Guid has no matching button: ", commandGuid)
		return
	}

	for _, v := range button.Commands {

		var err error
		switch vcasted := v.(type) {
		case *config.ShortcutCommand:
			err = executeKeyboardShortcut(vcasted)
		case *config.ScriptCommand:
			err = executeScript(vcasted)
		case *config.ShellCommand:
			err = executeShellCommand(vcasted)
		default:
			log.Println("Unsupported command type specified")
		}

		if err != nil {
			log.Printf("Error executing command \"%s\" for button \"%s\"\n", v.GetCommandType(), button.Name)
			return
		}
	}
}

func executeKeyboardShortcut(cmd *config.ShortcutCommand) error {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		return err
	}

	// According to https://github.com/micmonay/keybd_event,
	// linux needs a 2 second delay.
	// @ToaruBaka - You *can* (and should) wait for the device to show up in /sys,
	//	but it's way easier to just wait for a couple seconds right after setting it up
	if runtime.GOOS == "linux" {
		time.Sleep(2 * time.Second)
	}

	/*

	 */

	keysStr := cmd.Keys
	individualKeys := strings.Split(keysStr, "+")

	kb.Clear()
	for _, key := range individualKeys {
		keyTrimmed := strings.Trim(key, " \t\n\r")
		keyTrimmed = strings.ToUpper(keyTrimmed)

		if keyTrimmed == "CTRL" {
			kb.HasCTRL(true)
		} else if keyTrimmed == "ALT" {
			kb.HasALT(true)
		} else if keyTrimmed == "CMD" {
			kb.HasSuper(true)
		} else if keyTrimmed == "SHIFT" {
			kb.HasSHIFT(true)
		} else {
			mapping, ok := keyMappings[keyTrimmed]

			if !ok {
				return errors.New(fmt.Sprintf("Key string \"%s\" is invalid. No mapping found.\n", keyTrimmed))
			}

			kb.SetKeys(mapping)
		}
	}

	err = kb.Press()
	if err != nil {
		return err
	}

	time.Sleep(100 * time.Millisecond)

	err = kb.Release()
	if err != nil {
		return err
	}

	return nil
}

func executeScript(cmd *config.ScriptCommand) error {
	return nil
}

func executeShellCommand(cmd *config.ShellCommand) error {
	return nil
}
