package commands

import (
	"bufio"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/rahulCoder9417/Redis-in-go/server/config"
	"github.com/rahulCoder9417/Redis-in-go/server/types"
)

func ReplConf(client *types.Client, parts []string) string {
	if len(parts) >= 3 && strings.ToUpper(parts[1]) == "ACK" {
		offset, err := strconv.ParseInt(parts[2], 10, 64)
		if err == nil {
			replica := types.FindReplica(client.Conn)
			if replica != nil {
				replica.AckOffset = offset
			}
		}
		return ""
	}
	return RespSimpleString("OK")
}

func PSync(conn net.Conn, parts []string) string {
	if len(parts) != 3 {
		return RespError("wrong number of arguments for 'PSYNC'")
	}

	response :=
		"+FULLRESYNC " +
			config.ServerConfig.ReplicationID +
			" " +
			strconv.FormatInt(config.ServerConfig.ReplicationOffset, 10) +
			"\r\n"

	if _, err := conn.Write([]byte(response)); err != nil {
		return ""
	}

	SendEmptyRDB(conn)
	return ""
}

func SendEmptyRDB(conn net.Conn) error {
	emptyRDB := []byte{
		0x52, 0x45, 0x44, 0x49, 0x53, 0x30, 0x30, 0x31, 0x31,
	}

	header := "$" + strconv.Itoa(len(emptyRDB)) + "\r\n"

	if _, err := conn.Write([]byte(header)); err != nil {
		return err
	}
	if _, err := conn.Write(emptyRDB); err != nil {
		return err
	}
	_, err := conn.Write([]byte("\r\n"))
	return err
}

func DiscardRDB(reader *bufio.Reader) error {
	line, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	line = strings.TrimSpace(line)
	size, err := strconv.Atoi(line[1:])
	if err != nil {
		return err
	}

	buf := make([]byte, size)
	if _, err = io.ReadFull(reader, buf); err != nil {
		return err
	}

	_, err = reader.ReadString('\n')
	return err
}
