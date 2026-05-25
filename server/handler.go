package server

import (
	"strconv"
	"strings"
)

func HandleCommand(message string) string {

	parts := strings.Split(message, " ")

	if len(parts) == 0 {
		return "-ERR empty command\r\n"
	}

	command := strings.ToUpper(parts[0])

	switch command {

	case "PING":
		return "+PONG\r\n"

	case "ECHO":

		if len(parts) < 2 {
			return "-ERR wrong number of arguments for 'ECHO'\r\n"
		}

		msg := strings.Join(parts[1:], " ")

		return "$" +
			strconv.Itoa(len(msg)) +
			"\r\n" +
			msg +
			"\r\n"

	default:
		return "-ERR unknown command\r\n"
	}
}