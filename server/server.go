package server

import (
	"fmt"
	"net"
	"github.com/rahulCoder9417/Redis-in-go/server/config"
)

func Start() {

	ln, err := net.Listen("tcp", ":" + config.ServerConfig.Port)
	if err != nil {
		panic(err)
	}

	fmt.Println("Redis server started on port " + config.ServerConfig.Port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		fmt.Println("New client connected:", conn.RemoteAddr())

		go HandleClient(conn)
	}
}