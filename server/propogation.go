package server

import "fmt"

func Propogate(parts []string){
	ReplicaMu.Lock()
	defer ReplicaMu.Unlock()

	payload  := EncodeRESP(parts)

	activeReplicas := []*Replica{}
	
	for _,replica := range Replicas{
		_, err :=
			replica.Conn.Write(
				[]byte(payload),
			)
		
		fmt.Println(
			"propagating:",
			parts,
		)
		if err != nil {

			fmt.Println(
				"replica disconnected:",
				err,
			)

			replica.Conn.Close()

			continue
		}

		activeReplicas =
			append(
				activeReplicas,
				replica,
			)
	}

	Replicas =activeReplicas
}