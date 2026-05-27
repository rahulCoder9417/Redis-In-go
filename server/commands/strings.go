package commands

import (
	"strconv"
	"strings"
	"time"
)

func Ping(parts []string) string {
	return RespSimpleString("PONG")
}

func Echo(parts []string) string {
	if len(parts) < 2 {
		return RespError("wrong number of arguments for 'ECHO'")
	}
	msg := strings.Join(parts[1:], " ")
	return RespBulkString(msg)
}

func Set(parts []string) string {
	if len(parts) < 3 {
		return RespError("wrong number of arguments for 'SET'")
	}

	key := parts[1]
	value := parts[2]

	var expiry time.Time

	if len(parts) >= 5 {
		if strings.ToUpper(parts[3]) == "EX" {
			seconds, err := strconv.Atoi(parts[4])
			if err != nil {
				return RespError("invalid expire time")
			}
			expiry = time.Now().Add(time.Duration(seconds) * time.Second)
		}
	}

	Mu.Lock()
	Store[key] = Value{
		Type:      "string",
		String:    value,
		ExpiresAt: expiry,
	}
	Mu.Unlock()

	return RespSimpleString("OK")
}

func Get(parts []string) string {
	if len(parts) < 2 {
		return RespError("wrong number of arguments for 'GET'")
	}

	key := parts[1]

	Mu.RLock()
	value, exists := Store[key]
	Mu.RUnlock()

	if !exists {
		return RespNull()
	}

	if value.Type != "string" {
		return RespError("WRONGTYPE Operation against wrong kind of value")
	}

	if !value.ExpiresAt.IsZero() && value.ExpiresAt.Before(time.Now()) {
		Mu.Lock()
		delete(Store, key)
		Mu.Unlock()
		return RespNull()
	}

	return RespBulkString(value.String)
}
