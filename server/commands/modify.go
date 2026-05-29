package commands

import (
	"strconv"
	"time"
)

func Incr(parts []string) string {

	if len(parts) != 2 {
		return RespError(
			"wrong number of arguments for 'INCR'",
		)
	}

	key := parts[1]

	Mu.Lock()
	defer Mu.Unlock()

	value, exists := Store[key]

	if !exists {

		Store[key] = Value{
			Type:   "string",
			String: "1",
		}

		return RespInteger(1)
	}

	if !value.ExpiresAt.IsZero() &&
		value.ExpiresAt.Before(time.Now()) {

		delete(Store, key)

		Store[key] = Value{
			Type:   "string",
			String: "1",
		}

		return RespInteger(1)
	}

	if value.Type != "string" {
		return RespError(
			"WRONGTYPE Operation against wrong kind of value",
		)
	}

	num, err := strconv.Atoi(
		value.String,
	)

	if err != nil {
		return RespError(
			"value is not an integer",
		)
	}

	num++

	value.String =
		strconv.Itoa(num)

	Store[key] = value

	return RespInteger(num)
}