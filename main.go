package main

import (
	"github.com/rahulCoder9417/Redis-in-go/server/config"
	"github.com/rahulCoder9417/Redis-in-go/server"
	"os"
)

func main() {

	config.InitConfig()

	args := os.Args

	for i := 0; i < len(args); i++ {

		if args[i] == "--port" && i+1 < len(args) {

			config.ServerConfig.Port =
				args[i+1]
		}

		if args[i] == "--replicaof" && i+2 < len(args) {

			config.ServerConfig.IsReplica = true

			config.ServerConfig.MasterHost =
				args[i+1]

			config.ServerConfig.MasterPort =
				args[i+2]
		}
	}


	server.Start()
}