package main

import (
	"github.com/rahulCoder9417/Redis-in-go/server/config"
	"github.com/rahulCoder9417/Redis-in-go/server"
	"os"
)

func main() {

	config.InitConfig()

	args := os.Args

	if len(args) >= 4 && args[1] == "--replicaof" {
		config.ServerConfig.IsReplica = true
		config.ServerConfig.MasterHost = args[2]
		config.ServerConfig.MasterPort = args[3]
	}

	server.Start()
}