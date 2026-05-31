package server

import (
	"fmt"
	"github.com/rahulCoder9417/Redis-in-go/server/config"
	"github.com/rahulCoder9417/Redis-in-go/server/types"
)

func Propogate(parts []string) {
	types.ReplicaMu.Lock()
	defer types.ReplicaMu.Unlock()

	payload := EncodeRESP(parts)

	config.ServerConfig.ReplicationOffset +=
		int64(len(payload))

	activeReplicas := []*types.Replica{}

	for _, replica := range types.Replicas {
		_, err := replica.Conn.Write([]byte(payload))

		fmt.Println("propagating:", parts)

		if err != nil {
			fmt.Println("replica disconnected:", err)
			replica.Conn.Close()
			continue
		}

		activeReplicas = append(activeReplicas, replica)
	}

	types.Replicas = activeReplicas
}
