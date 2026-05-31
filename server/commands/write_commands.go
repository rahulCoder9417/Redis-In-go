package commands
import "strings"


func IsWritingCommand(
	cmd string
)bool{
	cmd = strings.ToUpper(cmd)

	switch cmd{
	case "SET",
		"INCR",

		"LPUSH",
		"RPUSH",
		"LPOP",
		"RPOP",

		"XADD":

		return true
	}

	return false
}