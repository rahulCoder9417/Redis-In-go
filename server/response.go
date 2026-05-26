package server

import "strconv"

func RespInteger(n int)string{
	return ":" + strconv.Itoa(n) + "\r\n"
}