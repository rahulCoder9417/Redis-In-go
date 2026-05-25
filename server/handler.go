package server

func HandleCommand(command string) string {

	switch command {

	case "PING":
		return "+PONG\r\n"

	default:
		return "-ERR unknown command\r\n"
	}
}