package streamdeck_server

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

	// TODO: Implement

}
