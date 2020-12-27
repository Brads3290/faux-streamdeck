package streamdeck_server

import (
	"faux-streamdeck/streamdeck_server/config"
	"github.com/micmonay/keybd_event"
	"log"
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
		case config.ShortcutCommand:
			err = executeKeyboardShortcut(vcasted)
		case config.ScriptCommand:
			err = executeScript(vcasted)
		case config.ShellCommand:
			err = executeShellCommand(vcasted)
		}

		if err != nil {
			log.Printf("Error executing command \"%s\" for button \"%s\"\n", v.GetCommandType(), button.Name)
			return
		}
	}
}

func executeKeyboardShortcut(cmd config.ShortcutCommand) error {

}

func executeScript(cmd config.ScriptCommand) error {

}

func executeShellCommand(cmd config.ShellCommand) error {

}
