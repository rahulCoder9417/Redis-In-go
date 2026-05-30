package commands

import (
	"strconv"
	"strings"

	"github.com/rahulCoder9417/Redis-in-go/server/config"
)

func Info(parts []string) string {

	if len(parts) != 2 {

		return RespError(
			"wrong number of arguments for 'INFO'",
		)
	}

	section :=
		strings.ToLower(parts[1])

	if section != "replication" {

		return RespBulkString("")
	}

	var response string

	if config.ServerConfig.IsReplica {

		response =
			"role:slave\r\n"

	} else {

		response =
			"role:master\r\n"

		response +=
			"master_replid:" +
				config.ServerConfig.ReplicationID +
				"\r\n"

		response +=
			"master_repl_offset:" +
				strconv.FormatInt(
					config.ServerConfig.ReplicationOffset,
					10,
				) +
				"\r\n"
	}

	return RespBulkString(
		response,
	)
}