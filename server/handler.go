package server

import (
	"strings"
	"github.com/rahulCoder9417/Redis-in-go/server/types"
	"github.com/rahulCoder9417/Redis-in-go/server/commands"
)

func HandleCommand(parts []string, client *types.Client) string {

	if len(parts) == 0 {
		return commands.RespError("empty command")
	}

	command := strings.ToUpper(parts[0])

	if client.InTransaction {

		switch command {

		case "EXEC":
			return commands.Exec(
				client,
				ExecuteImmediate,
			)
		case "DISCARD":
			return commands.Discard(client)
		case "MULTI":
			return commands.RespError("MULTI calls can not be nested")

		default:

			if !IsValidCommand(command) {

				client.HasTransactionError = true

				return commands.RespError(
					"unknown command",
				)
			}

			client.QueuedCommands =
				append(
					client.QueuedCommands,
					parts,
				)

			return commands.RespSimpleString(
					"QUEUED",
			)
		}
	}
	
	return ExecuteImmediate(client, parts)
}


func ExecuteImmediate(client *types.Client, parts []string) string {
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
	case "INCR":
		return commands.Incr(parts)
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
	case "MULTI":
		return commands.Multi(client)
	case "XADD":
		return commands.XAdd(parts)
	case "XRANGE":
		return commands.XRange(parts)
	case "XREAD":
		return commands.XRead(parts)
	case "DISCARD":
		return commands.Discard(client)
	case "EXEC":
		return commands.Exec(
			client,
			ExecuteImmediate,
		)
	case "WATCH":
		return commands.Watch(
			client,
			parts,
		)
	case "UNWATCH":
		return commands.UnWatch(client)
	case "INFO":
		return commands.Info(parts)
	case "REPLCONF":
		return commands.ReplConf(parts)

	case "PSYNC":
		return commands.PSync(parts)
	default:
		return commands.RespError("unknown command")
	}
}

func IsValidCommand(
	command string,
) bool {

	switch command {

	case "PING",
		"LRANGE",
		"ECHO",
		"SET",
		"GET",
		"INCR",
		"RPUSH",
		"LPUSH",
		"BLPOP",
		"LLEN",
		"LINDEX",
		"LPOP",
		"RPOP",
		"TYPE",
		"MULTI",
		"EXEC",
		"DISCARD",
		"XADD",
		"XRANGE",
		"XREAD",
		"WATCH",
		"UNWATCH",
		"INFO",
		"REPLCONF",
		"PSYNC":
		return true
	}

	return false
}