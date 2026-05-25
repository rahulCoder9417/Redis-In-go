package server

import (
	"strconv"
	"strings"
	"time"
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
    

	case "SET":
		if len(parts)<3{
			return "-ERR wrong number of arguments for 'SET'\r\n"
		}

		key := parts[1]
		value := parts[2]
		var expiry time.time

if len(parts)>=5{
			if strings.ToUpper(parts[3]) == "EX" {
				seconds,err := strconv.Atoi(parts[4])
				if err != nil {
					return "-ERR invalid expire time\r\n"
				}

				expiry = time.Now().Add(time.Duration(seconds) * time.Second)
			}
		}


		Mu.Lock()
		defer Mu.Unlock()

		
		Store[key] = Value{Data: value, ExpiresAt: expiry}

		return "+OK\r\n"


	case "GET":
		if len(parts)<2{
			return "-ERR wrong number of arguments for 'GET'\r\n"
		}
		
		key := parts[1]
		Mu.RLock()
		defer Mu.RUnlock()
		value, exists := Store[key]
		if !exists {
			return "$-1\r\n"
		}
		
		return "$" +
			strconv.Itoa(len(value.Data)) +
			"\r\n" +
			value.Data +
			"\r\n"
	default:
		return "-ERR unknown command\r\n"
	}
}