package types
import "net"

type Client struct {
	Conn net.Conn
	InTransaction bool
	QueuedCommands [][]string
	HasTransactionError bool
}