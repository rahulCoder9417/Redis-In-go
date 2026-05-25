package server

import (
	"sync"
	"time"
)

type Value struct{
	Data string
	ExpiresAt time.Time
}


var(

Store = map[string]Value{}
Mu sync.RWMutex
)