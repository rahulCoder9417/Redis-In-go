package commands

import (
	"sync"
	"time"
)

type StreamEntry struct{
	ID string
	Fields map[string]string
}

type Value struct {
	Type string

	String string
	List   []string
	Stream []StreamEntry

	ExpiresAt time.Time
}

var (
	Store         = map[string]Value{}
	Mu            sync.RWMutex
	ListWaiters   = map[string][]chan string{}
	StreamWaiters = map[string][]chan StreamEntry{}
)

func RemoveWaiter(key string, target chan string) {
	Mu.Lock()
	defer Mu.Unlock()
	waiters := ListWaiters[key]
	for i, waiter := range waiters {
		if waiter == target {
			ListWaiters[key] = append(waiters[:i], waiters[i+1:]...)
			break
		}
	}
}


func RemoveStreamWaiter(
	key string,
	target chan StreamEntry,
) {
	Mu.Lock()
	defer Mu.Unlock()
	waiters := StreamWaiters[key]
	for i, waiter := range waiters {
		if waiter == target {
			StreamWaiters[key] = append(waiters[:i], waiters[i+1:]...)
			break
		}
	}
}