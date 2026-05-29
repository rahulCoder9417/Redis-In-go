package commands


import "github.com/rahulCoder9417/Redis-in-go/server/types"

func Multi(client *types.Client)string {

	if client.InTransaction {
		return RespError("ERR MULTI calls can not be nested")
	}

	client.InTransaction = true
	
	client.QueuedCommands = nil
	return RespSimpleString("OK")
}

func Exec(
	client *types.Client,
	executor func(
		*types.Client,
		[]string,
	) string,
) string {

	if !client.InTransaction {

		return RespError(
			"EXEC without MULTI",
		)
	}

	queued := client.QueuedCommands

	client.InTransaction = false
	client.QueuedCommands = nil

	if len(queued) == 0 {
		return RespSimpleString("*0")
	}


	responses := []string{}

	for _, cmd := range queued {

		resp :=
			executor(
				client,
				cmd,
			)

		responses =
			append(
				responses,
				resp,
			)
	}

	return RespRawArray(
		responses,
	)
}