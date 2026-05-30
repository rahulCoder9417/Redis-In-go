package server

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/rahulCoder9417/Redis-in-go/server/config"
)

func ConnectToMaster() {

	address :=
		config.ServerConfig.MasterHost +
			":" +
			config.ServerConfig.MasterPort

	conn, err := net.Dial(
		"tcp",
		address,
	)

	if err != nil {

		fmt.Println(
			"failed connecting master:",
			err,
		)

		return
	}

	fmt.Println(
		"connected to master:",
		address,
	)

	reader :=
		bufio.NewReader(conn)

	// 1 PING

	SendCommand(
		conn,
		reader,
		[]string{
			"PING",
		},
	)

	// 2 REPLCONF listening-port

	SendCommand(
		conn,
		reader,
		[]string{
			"REPLCONF",
			"listening-port",
			config.ServerConfig.Port,
		},
	)

	// 3 REPLCONF capa psync2

	SendCommand(
		conn,
		reader,
		[]string{
			"REPLCONF",
			"capa",
			"psync2",
		},
	)

	// 4 PSYNC

	conn.Write(
		[]byte(
			EncodeRESP(
				[]string{
					"PSYNC",
					"?",
					"-1",
				},
			),
		),
	)

	err =HandleFullResync(
			reader,
		)

	if err != nil {

		fmt.Println(
			"fullresync failed:",
			err,
		)

		return
	}

	fmt.Println(
		"stored replication id:",
		config.ServerConfig.ReplicationID,
	)

	fmt.Println(
		"stored offset:",
		config.ServerConfig.ReplicationOffset,
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

func HandleFullResync(
	reader *bufio.Reader,
) error {

	line, err :=
		reader.ReadString('\n')

	if err != nil {

		return err
	}

	line =
		strings.TrimSpace(
			line,
		)

	fmt.Println(
		"master:",
		line,
	)

	parts :=
		strings.Split(
			line,
			" ",
		)

	if len(parts) != 3 {

		return fmt.Errorf(
			"invalid FULLRESYNC response",
		)
	}

	if parts[0] != "+FULLRESYNC" {

		return fmt.Errorf(
			"expected FULLRESYNC",
		)
	}

	config.ServerConfig.ReplicationID =
		parts[1]

	offset, err :=
		strconv.ParseInt(
			parts[2],
			10,
			64,
		)

	if err != nil {

		return err
	}

	config.ServerConfig.ReplicationOffset =
		offset

	return nil
}