package main

import (
	"github.com/rahulCoder9417/Redis-in-go/server"
	"os"
)

func main() {

	server.InitConfig()

	args := os.Args

	if len(args) >= 4 && args[1] == "--replicaof" {
		server.ServerConfig.IsReplica = true
		server.ServerConfig.MasterHost = args[2]
		server.ServerConfig.MasterPort = args[3]
	}

	server.Start()
}