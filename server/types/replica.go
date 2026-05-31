package types

import (
	"net"
	"sync"
)

type Replica struct {
	Conn      net.Conn
	Offset    int64
	AckOffset int64
}

var (
	Replicas  []*Replica
	ReplicaMu sync.Mutex
)

func AddReplica(conn net.Conn) {
	ReplicaMu.Lock()
	defer ReplicaMu.Unlock()

	Replicas = append(Replicas, &Replica{Conn: conn})
}

func RemoveReplica(conn net.Conn) {
	ReplicaMu.Lock()
	defer ReplicaMu.Unlock()

	filtered := []*Replica{}
	for _, replica := range Replicas {
		if replica.Conn != conn {
			filtered = append(filtered, replica)
		}
	}
	Replicas = filtered
}

func FindReplica(conn net.Conn) *Replica {
	ReplicaMu.Lock()
	defer ReplicaMu.Unlock()

	for _, replica := range Replicas {
		if replica.Conn == conn {
			return replica
		}
	}
	return nil
}
