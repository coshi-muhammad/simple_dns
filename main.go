package main

import (
	"fmt"
	"log"
	"net"
	"net/netip"
)

// constants
const MAX_BUFFER_SIZE = 4096

// types
type NameSpace struct {
	label      string
	address    netip.Addr
	sub_spaces []*NameSpace
}

func main() {
	// 1. Resolve the UDP address (specify protocol and port)
	addr, err := net.ResolveUDPAddr("udp", ":8080")
	if err != nil {
		log.Fatal("Could not resolve address:", err)
	}

	// 2. Start listening for UDP packets on the resolved address
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal("Listen failed:", err)
	}
	defer conn.Close()

	fmt.Println("DNS server listening on :8080")

	// 3. Create a buffer to hold incoming data
	buffer := make([]byte, MAX_BUFFER_SIZE)

	for {
		// 4. Read data from the connection
		// ReadFromUDP returns the number of bytes read and the sender's address
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("Read error: %v\n", err)
			continue
		}

		// 5. Process the received message
		message := string(buffer[:n])
		fmt.Printf("Received %d bytes from %s: %s\n", n, clientAddr, message)

		// Optional: Echo a response back to the client
		// _, err = conn.WriteToUDP([]byte("Message received"), clientAddr)
	}
}
