package commands

import (
	"strconv"
	"time"
)

func RPush(parts []string) string {
	if len(parts) < 3 {
		return RespError("wrong number of arguments for 'RPUSH'")
	}

	key := parts[1]
	values := parts[2:]

	Mu.Lock()
	defer Mu.Unlock()

	v, exists := Store[key]

	if !exists {
		v = Value{
			Type: "list",
			List: []string{},
		}
	}

	if v.Type != "list" {
		return RespError("WRONGTYPE Operation against wrong kind of value")
	}

	for _, value := range values {
		waiters := ListWaiters[key]
		if len(waiters) > 0 {
			waiter := waiters[0]
			ListWaiters[key] = waiters[1:]
			waiter <- value
		} else {
			v.List = append(v.List, value)
		}
	}

	Store[key] = v

	return RespInteger(len(v.List))
}

func LPush(parts []string) string {
	if len(parts) < 3 {
		return RespError("wrong number of arguments for 'LPUSH'")
	}

	key := parts[1]
	values := parts[2:]

	Mu.Lock()
	defer Mu.Unlock()

	v, exists := Store[key]

	if !exists {
		v = Value{
			Type: "list",
			List: []string{},
		}
	}

	if v.Type != "list" {
		return RespError("WRONGTYPE Operation against wrong kind of value")
	}

	for _, value := range values {
		waiters := ListWaiters[key]
		if len(waiters) > 0 {
			waiter := waiters[0]
			ListWaiters[key] = waiters[1:]
			waiter <- value
		} else {
			v.List = append([]string{value}, v.List...)
		}
	}

	Store[key] = v

	return RespInteger(len(v.List))
}

func BLPop(parts []string) string {
	if len(parts) < 3 {
		return RespError("wrong number of arguments for 'BLPOP'")
	}

	key := parts[1]

	timeout, err := strconv.Atoi(parts[2])
	if err != nil {
		return RespError("invalid timeout")
	}

	Mu.Lock()

	v, exists := Store[key]

	if exists {
		if v.Type != "list" {
			Mu.Unlock()
			return RespError("WRONGTYPE Operation against wrong kind of value")
		}

		if len(v.List) > 0 {
			item := v.List[0]
			v.List = v.List[1:]

			if len(v.List) == 0 {
				delete(Store, key)
			} else {
				Store[key] = v
			}

			Mu.Unlock()
			return RespBulkString(item)
		}
	}

	ch := make(chan string)
	ListWaiters[key] = append(ListWaiters[key], ch)
	Mu.Unlock()

	if timeout == 0 {
		item := <-ch
		return RespBulkString(item)
	}

	select {
	case item := <-ch:
		return RespBulkString(item)
	case <-time.After(time.Duration(timeout) * time.Second):
		return RespNull()
	}
}

func LLen(parts []string) string {
	if len(parts) < 2 {
		return RespError("wrong number of arguments for 'LLEN'")
	}

	key := parts[1]

	Mu.RLock()
	value, exists := Store[key]
	Mu.RUnlock()

	if !exists {
		return RespInteger(0)
	}

	if value.Type != "list" {
		return RespError("WRONGTYPE Operation against wrong kind of value")
	}

	return RespInteger(len(value.List))
}

func LIndex(parts []string) string {
	if len(parts) < 3 {
		return RespError("wrong number of arguments for 'LINDEX'")
	}

	key := parts[1]

	index, err := strconv.Atoi(parts[2])
	if err != nil {
		return RespError("index out of range")
	}

	Mu.RLock()
	value, exists := Store[key]
	Mu.RUnlock()

	if !exists {
		return RespNull()
	}

	if value.Type != "list" {
		return RespError("WRONGTYPE Operation against wrong kind of value")
	}

	if index < 0 {
		index = len(value.List) + index
	}

	if index < 0 || index >= len(value.List) {
		return RespNull()
	}

	return RespBulkString(value.List[index])
}

func LPop(parts []string) string {
	if len(parts) < 2 {
		return RespError("wrong number of arguments for 'LPOP'")
	}

	key := parts[1]

	Mu.Lock()
	defer Mu.Unlock()

	v, exists := Store[key]

	if !exists {
		return RespNull()
	}

	if v.Type != "list" {
		return RespError("WRONGTYPE Operation against wrong kind of value")
	}

	if len(v.List) == 0 {
		return RespNull()
	}

	item := v.List[0]
	v.List = v.List[1:]

	if len(v.List) == 0 {
		delete(Store, key)
	} else {
		Store[key] = v
	}

	return RespBulkString(item)
}

func RPop(parts []string) string {
	if len(parts) < 2 {
		return RespError("wrong number of arguments for 'RPOP'")
	}

	key := parts[1]

	Mu.Lock()
	defer Mu.Unlock()

	v, exists := Store[key]

	if !exists {
		return RespNull()
	}

	if v.Type != "list" {
		return RespError("WRONGTYPE Operation against wrong kind of value")
	}

	if len(v.List) == 0 {
		return RespNull()
	}

	last := v.List[len(v.List)-1]
	v.List = v.List[:len(v.List)-1]

	if len(v.List) == 0 {
		delete(Store, key)
	} else {
		Store[key] = v
	}

	return RespBulkString(last)
}
