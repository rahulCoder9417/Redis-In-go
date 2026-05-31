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

func PSync(conn net.Conn, parts []string)string{
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
// it the replica mis-reads the first propagated write as the RDB header.
	if _, err := conn.Write([]byte(response)); err != nil {
		return ""
	}

	SendEmptyRDB(conn)

	// everything was written directly to the connection.
	return ""
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
	if err != nil {
		return err
	}

	_, err =
		conn.Write(
			[]byte("\r\n"),
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

	_, err = reader.ReadString('\n')

	if err != nil {
		return err
	}

	return nil
}