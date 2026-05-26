package server

import (
	"sync"
	"time"
)

type Value struct{
	Type string

	String string
	List []string

	ExpiresAt time.Time
}


var(
	Store = map[string]Value{}
	Mu sync.RWMutex
)
