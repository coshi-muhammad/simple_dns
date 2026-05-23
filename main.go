package main

import (
	"fmt"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error Listen:", err.Error())
	}
	defer listener.Close()
	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("Error Connection:", err.Error())
			continue
		}
		buff := make([]byte, 1024)
		length, err := connection.Read(buff)
		if err != nil {
			fmt.Println("Error Reading Connection:", err.Error())
			continue
		}
		fmt.Printf("Reseved:%s\n", string(buff[:length]))
		connection.Write([]byte("Message resived."))
		connection.Close()
	}
}
