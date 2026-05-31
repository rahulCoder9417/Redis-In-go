package commands

import (
	"strconv"
	"time"

	"github.com/rahulCoder9417/Redis-in-go/server/types"
	"github.com/rahulCoder9417/Redis-in-go/server/config"
)

func Wait(
	parts []string,
)string{
	if len(parts)!=3{
		return RespError(
			"wrong numbers of arguments for WAIT",
		)
	}

	required, err :=
		strconv.Atoi(
			parts[1],
		)

	if err != nil {

		return RespError(
			"invalid replica count",
		)
	}

	timeoutMs, err :=
		strconv.Atoi(
			parts[2],
		)

	if err != nil {

		return RespError(
			"invalid timeout",
		)
	}

	targetOffset :=
		config.ServerConfig.
			ReplicationOffset

	deadline :=
		time.Now().Add(
			time.Duration(
				timeoutMs,
			) *
				time.Millisecond,
		)

	for {

		count := 0

		types.ReplicaMu.Lock()

		for _, replica := range types.Replicas {

			if replica.AckOffset >=
				targetOffset {

				count++
			}
		}

		types.ReplicaMu.Unlock()

		if count >= required {

			return RespInteger(
				count,
			)
		}

		if time.Now().After(
			deadline,
		) {

			return RespInteger(
				count,
			)
		}

		time.Sleep(
			10 * time.Millisecond,
		)
	}
}