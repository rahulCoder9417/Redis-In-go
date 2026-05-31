package server

import (
	"bufio"
	"fmt"
	"net"
	"github.com/rahulCoder9417/Redis-in-go/server/types"
)

func HandleClient(conn net.Conn) {
	

	client := &types.Client{
		Conn: conn,
	}
	defer func() {
		client.WatchedKeys = nil
		types.RemoveReplica(conn)
		conn.Close()
	}()
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