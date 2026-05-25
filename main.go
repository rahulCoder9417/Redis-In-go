package main

import (
	"fmt"
	"net"
)

func main(){
	ln,err := net.Listen("tcp",":6379")
	if err!=nil{
		panic(err)
	}

	fmt.Println("Server started on :6379")

	for{
		conn,err := ln.Accept()
		if err!=nil{
			fmt.Println("Error accepting Connection",err)
			continue
		}
		fmt.Println("Client Connected",conn.RemoteAddr())
	}
}
