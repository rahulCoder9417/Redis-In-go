package server

import (
	"bufio"
	"fmt"
	"net"
	"github.com/rahulCoder9417/Redis-in-go/server/config"
)

func ConnectToMaster(){
	address:=
		config.ServerConfig.MasterHost +
			":" + 
			config.ServerConfig.MasterPort
	
	conn, err := net.Dial(
		"tcp",
		address,
	)

	if err != nil{
		fmt.Println(
			"failed Connecting master:", err,
		)
		return
	}
	fmt.Println(
		"connected to master:",
		address,
	)

	reader :=
		bufio.NewReader(conn)

	SendCommand(
		conn,
		reader,
		[]string{"PING"},
	)

	SendCommand(
		conn,
		reader,
		[]string{
			"REPLCONF",
			"listening-port",
			config.ServerConfig.Port,
		},
	)

	SendCommand(
		conn,
		reader,
		[]string{
			"REPLCONF",
			"capa",
			"psync2",
		},
	)

	SendCommand(
		conn,
		reader,
		[]string{
			"PSYNC",
			"?",
			"-1",
		},
	)

	fmt.Println(
		"replica handshake completed",
	)
}

func SendCommand(
	conn net.Conn,
	reader *bufio.Reader,
	cmd []string,
) {

	encoded :=
		EncodeRESP(cmd)

	conn.Write(
		[]byte(encoded),
	)

	response, err :=
		reader.ReadString('\n')

	if err != nil {

		fmt.Println(
			"handshake failed:",
			err,
		)

		return
	}

	fmt.Print(
		"master:",
		response,
	)
}