package server

import (
	"strings"
	"github.com/rahulCoder9417/Redis-in-go/server/types"
	"github.com/rahulCoder9417/Redis-in-go/server/commands"
	"github.com/rahulCoder9417/Redis-in-go/server/config"
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

	var response string

	switch command {
	case "PING":
		response = commands.Ping(parts)
	case "LRANGE":
		response = commands.LRange(parts)
	case "ECHO":
		response = commands.Echo(parts)
	case "SET":
		response = commands.Set(parts)
	case "GET":
		response = commands.Get(parts)
	case "INCR":
		response = commands.Incr(parts)
	case "RPUSH":
		response = commands.RPush(parts)
	case "LPUSH":
		response = commands.LPush(parts)
	case "BLPOP":
		response = commands.BLPop(parts)
	case "LLEN":
		response = commands.LLen(parts)
	case "LINDEX":
		response = commands.LIndex(parts)
	case "LPOP":
		response = commands.LPop(parts)
	case "RPOP":
		response = commands.RPop(parts)
	case "TYPE":
		response = commands.Type(parts)
	case "MULTI":
		response = commands.Multi(client)
	case "XADD":
		response = commands.XAdd(parts)
	case "XRANGE":
		response = commands.XRange(parts)
	case "XREAD":
		response = commands.XRead(parts)
	case "DISCARD":
		response = commands.Discard(client)
	case "EXEC":
		response = commands.Exec(
			client,
			ExecuteImmediate,
		)
	case "WATCH":
		response = commands.Watch(
			client,
			parts,
		)
	case "UNWATCH":
		response = commands.UnWatch(client)
	case "INFO":
		response = commands.Info(parts)
	case "REPLCONF":
		response = commands.ReplConf(client,parts)

	case "PSYNC":
		response = commands.PSync(client.Conn, parts)
		types.AddReplica(client.Conn)
	case "WAIT":
		response = commands.Wait(parts)
	default:
		return commands.RespError("unknown command")
	}



	if commands.IsWritingCommand(
		command,
	) && !config.ServerConfig.IsReplica {

		Propogate(parts)
	}

	return response
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
		"PSYNC",
		"WAIT":
		return true
	}

	return false
}