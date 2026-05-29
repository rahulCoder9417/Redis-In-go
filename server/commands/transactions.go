package commands


import "github.com/rahulCoder9417/Redis-in-go/server/types"

func Multi(client *types.Client)string {

	if client.InTransaction {
		return RespError(" MULTI calls can not be nested")
	}

	client.InTransaction = true
	client.HasTransactionError = false
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


	if client.HasTransactionError {

		client.InTransaction = false
		client.QueuedCommands = nil
		client.HasTransactionError = false

		return RespError(
			"EXECABORT Transaction discarded because of previous errors",
		)
	}

	queued := client.QueuedCommands

	client.InTransaction = false
	client.QueuedCommands = nil
	client.HasTransactionError = false
	if len(queued) == 0 {
		return "*0\r\n"
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

func Discard(client *types.Client) string {
	if !client.InTransaction {
		return RespError(" DISCARD without MULTI")
	}
	client.HasTransactionError = false
	
	client.InTransaction = false
	client.QueuedCommands = nil
	return RespSimpleString("OK")
}
