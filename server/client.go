package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func HandleClient(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {

		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Client disconnected:", conn.RemoteAddr())
			return
		}

		message = strings.TrimSpace(message)

		fmt.Printf("[%s] %s\n", conn.RemoteAddr(), message)

		response := HandleCommand(message)

		conn.Write([]byte(response))
	}
}