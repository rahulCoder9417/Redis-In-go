package server

import (
	"strconv"
	"strings"
	"time"
)

func HandleCommand(parts []string) string {

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
		var expiry time.Time

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

		
		Store[key] = Value{Type: "string",String: value, ExpiresAt: expiry}

		return "+OK\r\n"


	case "GET":
		if len(parts)<2{
			return "-ERR wrong number of arguments for 'GET'\r\n"
		}
		
		key := parts[1]
		Mu.RLock()
		value, exists := Store[key]
		Mu.RUnlock()
		if !exists {
			return "$-1\r\n"
		}

		if value.Type != "string" {
			return "-ERR value is not a string\r\n"
		}
		
		if !value.ExpiresAt.IsZero() && value.ExpiresAt.Before(time.Now()) {
			Mu.Lock()
			delete(Store, key)
			Mu.Unlock()
			return "$-1\r\n"
		}

		return "$" +
			strconv.Itoa(len(value.String)) +
			"\r\n" +
			value.String +
			"\r\n"
	
	
	case "RPUSH":
		if len(parts)<3{
			return "-ERR wrong number of arguments for 'RPUSH'\r\n"
		}
		
		key := parts[1]
		values:= parts[2:]


		Mu.Lock()
		defer Mu.Unlock()
		
		v,exists := Store[key]
		
		if !exists{
			v=Value{
				Type:"list",
				List: []string{},
			}
		}
		
		if v.Type != "list"{
			return "-WRONGTYPE Operation against wrong kind of value\r\n"
		}

		v.List = append(v.List, values...)
		Store[key] = v
		
		return RespInteger(len(v.List))

	case "LPUSH":
		if len(parts)<3{
			return "-ERR wrong number of arguments for 'LPUSH'\r\n"
		}

		key := parts[1]
		values:=parts[2:]

		Mu.Lock()
		defer Mu.Unlock()

		v,exists :=Store[key]

		if !exists{
			v = Value{
				Type:"list",
				List: []string{},
			}
		}
		
		if v.Type != "list"{
			return "-WRONGTYPE Operation against wrong kind of value\r\n"
		}

		for _,value:=range values{
			v.List = append([]string{value}, v.List...)
		}

		Store[key] = v
		
		return RespInteger(len(v.List))
	default:
		return "-ERR unknown command\r\n"
	}
}