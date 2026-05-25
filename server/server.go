package server

import (
	"fmt"
	"net"
)

func Start() {

	ln, err := net.Listen("tcp", ":6379")
	if err != nil {
		panic(err)
	}

	fmt.Println("Redis server started on port 6379")

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