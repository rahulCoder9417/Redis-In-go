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


func Watch(
	client *types.Client,
	parts []string,
) string {
	if len(parts)<2{
		return RespError("wrong number of arguments for 'WATCH' command")
	}
	
	if client.WatchedKeys == nil {
		client.WatchedKeys = make(map[string]int64)
	}
	
	Mu.RLock()
	defer Mu.RUnlock()
	for _, key := range parts[1:] {
		version := KeyVersions[key]

		client.WatchedKeys[key] =version
	}
	
	return RespSimpleString("OK")
}


func UnWatch(client *types.Client) string {
	
	
	client.WatchedKeys = nil
	return RespSimpleString("OK")
}
