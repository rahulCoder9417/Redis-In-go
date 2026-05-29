package commands

import (
	"strconv"
)

func RespSimpleString(s string) string {
	return "+" + s + "\r\n"
}

func RespError(s string) string {
	return "-ERR " + s + "\r\n"
}

func RespBulkString(s string) string {
	return "$" +
		strconv.Itoa(len(s)) +
		"\r\n" +
		s +
		"\r\n"
}

func RespNull() string {
	return "$-1\r\n"
}

func RespInteger(n int) string {
	return ":" + strconv.Itoa(n) + "\r\n"
}

func RespArray(arr []string) string {

	resp := "*" + strconv.Itoa(len(arr)) + "\r\n"

	for _, item := range arr {
		resp += RespBulkString(item)
	}

	return resp
}

func RespRawArray(
	responses []string,
) string {

	resp :=
		"*" +
		strconv.Itoa(
			len(responses),
		) +
		"\r\n"

	for _, r := range responses {

		resp += r
	}

	return resp
}