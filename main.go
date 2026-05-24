package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"net/netip"
)

// constants
const MAX_BUFFER_SIZE = 4096

// enums
type OpCode int

const (
	QUERY OpCode = iota
	IQUERY
	STATUS
	RESERVED
	NOTIFY
	UPDATE
)

type RCode int

const (
	NO_ERROR RCode = iota
	FORMAT_ERROR
	SERVER_ERROR
	NAME_ERROR
	NOT_IMPLEMENTED
	REFUSED

	//NOTE: added for completion posibly removed if not implemented
	XY_DOMAIN
	XY_RR_SET
	NX_RR_SET
	NOT_AUTH
	NOT_ZONE
	// added for completion posibly removed if not implemented
)

// types
type NameSpace struct {
	label      string
	address    netip.Addr
	sub_spaces []*NameSpace
} //used to store the database information corelating a domain with an ip address

type MessageHeader struct {
	id          uint16
	qr          bool
	flags_codes uint16
	op_code     OpCode
	aa          bool
	tc          bool
	rd          bool
	ra          bool
	r_code      RCode
	qd_count    uint16
	an_count    uint16
	ns_count    uint16
	ar_count    uint16
}

func (mh *MessageHeader) decode(header_raw []byte) error {
	if len(header_raw) != 12 {
		return fmt.Errorf("Invalid header size")
	}
	mh.id = binary.BigEndian.Uint16(header_raw[:2])
	mh.flags_codes = binary.BigEndian.Uint16(header_raw[2:4])
	mh.qr = ((mh.flags_codes & 0x8000) >> 15) == 1
	mh.op_code = OpCode((mh.flags_codes & 0x7800) >> 11)
	mh.aa = ((mh.flags_codes & 0x0400) >> 10) == 1
	mh.tc = ((mh.flags_codes & 0x0200) >> 9) == 1
	mh.rd = ((mh.flags_codes & 0x0100) >> 8) == 1
	mh.ra = ((mh.flags_codes & 0x0080) >> 7) == 1
	mh.r_code = RCode(mh.flags_codes & 0x000F)
	mh.qd_count = binary.BigEndian.Uint16(header_raw[4:6])
	mh.an_count = binary.BigEndian.Uint16(header_raw[6:8])
	mh.ns_count = binary.BigEndian.Uint16(header_raw[8:10])
	mh.ar_count = binary.BigEndian.Uint16(header_raw[10:])
	return nil
}
func (mh *MessageHeader) encode() []byte {
	header_raw := make([]byte, 12)
	binary.BigEndian.PutUint16(header_raw[:2], mh.id)
	binary.BigEndian.PutUint16(header_raw[2:4], mh.flags_codes)
	binary.BigEndian.PutUint16(header_raw[4:6], mh.qd_count)
	binary.BigEndian.PutUint16(header_raw[6:8], mh.an_count)
	binary.BigEndian.PutUint16(header_raw[8:10], mh.ns_count)
	binary.BigEndian.PutUint16(header_raw[10:], mh.ar_count)
	return header_raw
}
func (mh *MessageHeader) print() {
	fmt.Println("message reseved :")
	fmt.Println("id: ", mh.id)
	fmt.Println("is query: ", !mh.qr)
	fmt.Println("operatoin code: ", mh.op_code)
	fmt.Println("answer is autoritative: ", mh.aa)
	fmt.Println("truncation happend: ", mh.tc)
	fmt.Println("recursion desired: ", mh.rd)
	fmt.Println("recursion available: ", mh.ra)
	fmt.Println("response code: ", mh.r_code)
	fmt.Println("question count: ", mh.qd_count)
	fmt.Println("answer count: ", mh.an_count)
	fmt.Println("authority count: ", mh.ns_count)
	fmt.Println("aditional count: ", mh.ar_count)
	fmt.Println("header end")
}

type Message struct {
	header MessageHeader
	//TODO: add the other sections

} //used to represent the shape of a dns message when sending or reseaving

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
		header_raw := buffer[:12]
		fmt.Printf("header hex: %v", header_raw)
		header := MessageHeader{}
		header.decode(header_raw)
		header.print()
		second_header_raw := header.encode()
		fmt.Printf("header hex: %v", second_header_raw)
		fmt.Printf("Received %d bytes from %s: %s\n", n, clientAddr, message)

		// Optional: Echo a response back to the client
		// _, err = conn.WriteToUDP([]byte("Message received"), clientAddr)
	}
}
