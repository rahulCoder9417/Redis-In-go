package commands

import (
	"strconv"
	"github.com/rahulCoder9417/Redis-in-go/server/config"
)

func ReplConf(parts []string)string{
	return RespSimpleString("OK")
}

func PSync(parts []string)string{
	if len(parts) != 3 {
		return RespError(
			"wrong number of arguments for 'PSYNC'",
		)
	}
	response:=
		"+FULLRESYNC "+
			config.ServerConfig.ReplicationID +
			" " +
			strconv.FormatInt(
				config.ServerConfig.ReplicationOffset,
				10,
			) +
			"\r\n"

	return response
}