package commands

func Type(parts []string) string {
	if len(parts) != 2 {
		return RespError("wrong number of arguments for 'TYPE'")
	}

	key := parts[1]

	Mu.RLock()
	value, exists := Store[key]
	Mu.RUnlock()

	if !exists {
		return RespSimpleString("none")
	}

	return RespSimpleString(value.Type)
}