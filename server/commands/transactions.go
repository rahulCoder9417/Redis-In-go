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
