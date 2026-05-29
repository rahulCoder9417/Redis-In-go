package server

import (
	"bufio"
	"fmt"
	"net"
)

func HandleClient(conn net.Conn) {
	defer conn.Close()

	client := &Client{
		Conn: conn,
	}

	reader := bufio.NewReader(conn)
	resp:=NewResp(reader)

	for {

		command, err := resp.Read()
		if err != nil {
			fmt.Println("Client disconnected:", conn.RemoteAddr())
			return
		}

		fmt.Printf("[%s] %s\n", conn.RemoteAddr(), command)

		response := HandleCommand(
			command,
			client,
		)

		conn.Write([]byte(response))
	}
}