package main

import (
	"fmt"
	"net"
	"bufio"
	"strings"
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
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn){
	defer conn.Close()

	reader := bufio.NewReader(conn)
	
	for {
		message,err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading message",conn.RemoteAddr(),err)
			return
		}

		message = strings.TrimSpace(message)

		fmt.Println("Message received",message)

		if message == "PING" {
			conn.Write([]byte("+PONG\r\n"))
		}else {
			conn.Write([]byte("-ERR unknown command\r\n"))
		}
	}
}
