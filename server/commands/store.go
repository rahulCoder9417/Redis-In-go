package commands

import (
	"sync"
	"time"
)

type Value struct {
	Type string

	String string
	List   []string

	ExpiresAt time.Time
}

var (
	Store       = map[string]Value{}
	Mu          sync.RWMutex
	ListWaiters = map[string][]chan string{}
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
