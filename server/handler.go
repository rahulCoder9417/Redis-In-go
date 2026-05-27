package server

import (
	"strings"

	"github.com/rahulCoder9417/Redis-in-go/server/commands"
)

func HandleCommand(parts []string) string {

	if len(parts) == 0 {
		return commands.RespError("empty command")
	}

	command := strings.ToUpper(parts[0])

	switch command {
	case "PING":
		return commands.Ping(parts)
	case "LRANGE":
		return commands.LRange(parts)
	case "ECHO":
		return commands.Echo(parts)
	case "SET":
		return commands.Set(parts)
	case "GET":
		return commands.Get(parts)
	case "RPUSH":
		return commands.RPush(parts)
	case "LPUSH":
		return commands.LPush(parts)
	case "BLPOP":
		return commands.BLPop(parts)
	case "LLEN":
		return commands.LLen(parts)
	case "LINDEX":
		return commands.LIndex(parts)
	case "LPOP":
		return commands.LPop(parts)
	case "RPOP":
		return commands.RPop(parts)
	case "TYPE":
		return commands.Type(parts)
	case "XADD":
		return commands.XAdd(parts)
	default:
		return commands.RespError("unknown command")
	}
}
