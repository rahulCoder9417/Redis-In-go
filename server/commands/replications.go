package commands

import (
	"strconv"
	"github.com/rahulCoder9417/Redis-in-go/server/config"
	"bufio"
	"io"
	"strings"
	"net"
)

func ReplConf(parts []string)string{
	return RespSimpleString("OK")
}

func PSync(parts []string)string{
	if len(parts) != 3 {
		return RespError(
			"wrong number of arguments for 'PSYNC'",
		)
	}
	response:=
		"+FULLRESYNC "+
			config.ServerConfig.ReplicationID +
			" " +
			strconv.FormatInt(
				config.ServerConfig.ReplicationOffset,
				10,
			) +
			"\r\n"

	return response
}



func SendEmptyRDB(conn net.Conn)error {
	emptyRDB := []byte{
		0x52,
		0x45,
		0x44,
		0x49,
		0x53,
		0x30,
		0x30,
		0x31,
		0x31,
	}

	header :=
		"$" +
		strconv.Itoa(
			len(emptyRDB),
		) +
		"\r\n"

	_, err :=
		conn.Write(
			[]byte(header),
		)

	if err != nil {
		return err
	}

	_, err =
		conn.Write(
			emptyRDB,
		)

	return err
}

func DiscardRDB(
	reader *bufio.Reader,
) error {

	line, err :=
		reader.ReadString('\n')

	if err != nil {
		return err
	}

	line =
		strings.TrimSpace(line)

	size, err :=
		strconv.Atoi(
			line[1:],
		)

	if err != nil {
		return err
	}

	buf :=
		make(
			[]byte,
			size,
		)

	_, err = io.ReadFull(reader, buf)

	if err != nil {
		return err
	}

	reader.ReadString('\n')

	return nil
}