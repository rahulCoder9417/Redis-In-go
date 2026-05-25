package server

import "sync"

type Value struct{
	Data string
}


var(

Store = map[string]Value{}
Mu sync.RWMutex
)