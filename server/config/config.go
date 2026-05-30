package config

import (
	"crypto/rand"
	"encoding/hex"
)

type Config struct {
	Port string

	IsReplica bool

	MasterHost string
	MasterPort string

	ReplicationID string
	ReplicationOffset int64
}

var ServerConfig Config

func InitConfig() {

	ServerConfig = Config{
		Port: "6380",
		
		IsReplica: false,

		ReplicationID: RandomID(),

		ReplicationOffset: 0,
	}
}

func RandomID() string {

	bytes := make([]byte, 20)

	rand.Read(bytes)

	return hex.EncodeToString(bytes)
}