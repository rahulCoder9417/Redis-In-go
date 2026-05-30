package server

import (
	"net"
    "sync"
)

type Replica struct {
	Conn net.Conn

	Offset int64

	AckOffset int64
}

var(
	Replicas []*Replica
	
	ReplicaMu sync.Mutex
)