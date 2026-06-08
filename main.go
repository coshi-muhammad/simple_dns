package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"net/netip"
	"slices"
	"strings"
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

type RRType uint16

const (
	ADDRESS_IPV4      RRType = 1
	ADDRESS_IPV6      RRType = 28
	NAME_SERVER       RRType = 2
	CANONICAL_NAME    RRType = 5
	START_OF_AUTORITY RRType = 6
	POINTER           RRType = 12
	MAIL_EXCHANGE     RRType = 15
	TEXT              RRType = 16
)

// types
type NameSpace struct {
	domain     string
	address    netip.Addr
	subdomains []NameSpace
}

func loadNameSpace() NameSpace {
	//TODO: at somepoint make it do a parce to the master file and give you this
	// but for now im just gonna hard code some values
	name_space := NameSpace{}
	name_space.subdomains = make([]NameSpace, 2)
	name_space.subdomains[0] = NameSpace{
		domain: "com",
	}
	com_name_space := &name_space.subdomains[0]
	com_name_space.subdomains = make([]NameSpace, 3)
	temp := netip.AddrFrom4([4]byte{192, 168, 1, 4})
	com_name_space.subdomains[0] = NameSpace{
		domain:  "coshi",
		address: temp,
	}
	temp = netip.AddrFrom4([4]byte{192, 168, 1, 10})
	com_name_space.subdomains[1] = NameSpace{
		domain:  "younes",
		address: temp,
	}
	temp = netip.AddrFrom4([4]byte{145, 255, 100, 14})
	com_name_space.subdomains[2] = NameSpace{
		domain:  "google",
		address: temp,
	}

	name_space.subdomains[1] = NameSpace{
		domain: "org",
	}
	org_name_space := &name_space.subdomains[1]
	org_name_space.subdomains = make([]NameSpace, 2)
	temp = netip.AddrFrom16([16]byte{
		0x32, 0x44, 0x0, 0x0, 0x3, 0x1, 0x0, 0x0, 0x0, 0xf3, 0x09, 0x3c, 0xaa, 0x0, 0x0, 0x04,
	})
	org_name_space.subdomains[0] = NameSpace{
		domain:  "potato",
		address: temp,
	}
	temp = netip.AddrFrom16([16]byte{
		0x16, 0x0, 0xf7, 0x3f, 0x1d, 0xe1, 0xdd, 0x0, 0x0, 0xf3, 0x04, 0xab, 0xca, 0x0, 0x0, 0x09,
	})
	org_name_space.subdomains[1] = NameSpace{
		domain:  "historyfacts",
		address: temp,
	}

	return name_space
}

func getAddress(space NameSpace, name []string) (netip.Addr, error) {
	for _, domain := range slices.Backward(name) {
		found := false
		for _, subspace := range space.subdomains {
			if subspace.domain == domain {
				found = true
				space = subspace
				break
			}
		}
		if !found {
			return netip.Addr{}, fmt.Errorf("Address not found")
		}
	}
	return space.address, nil
}

type MessageHeader struct {
	id       uint16
	qr       bool
	op_code  OpCode
	aa       bool
	tc       bool
	rd       bool
	ra       bool
	r_code   RCode
	qd_count uint16
	an_count uint16
	ns_count uint16
	ar_count uint16
}

func decodeHeader(header_raw []byte) (MessageHeader, error) {
	if len(header_raw) != 12 {
		return MessageHeader{}, fmt.Errorf("Invalid header size")
	}
	mh := MessageHeader{}
	mh.id = binary.BigEndian.Uint16(header_raw[:2])
	flags_codes := binary.BigEndian.Uint16(header_raw[2:4])
	mh.qr = ((flags_codes & 0x8000) >> 15) == 1
	mh.op_code = OpCode((flags_codes & 0x7800) >> 11)
	mh.aa = ((flags_codes & 0x0400) >> 10) == 1
	mh.tc = ((flags_codes & 0x0200) >> 9) == 1
	mh.rd = ((flags_codes & 0x0100) >> 8) == 1
	mh.ra = ((flags_codes & 0x0080) >> 7) == 1
	mh.r_code = RCode(flags_codes & 0x000F)
	mh.qd_count = binary.BigEndian.Uint16(header_raw[4:6])
	mh.an_count = binary.BigEndian.Uint16(header_raw[6:8])
	mh.ns_count = binary.BigEndian.Uint16(header_raw[8:10])
	mh.ar_count = binary.BigEndian.Uint16(header_raw[10:12])
	return mh, nil
}
func encodeHeader(mh *MessageHeader) []byte {
	header_raw := make([]byte, 0, 12)
	header_raw = binary.BigEndian.AppendUint16(header_raw, mh.id)
	var flags_codes uint16
	if mh.qr {
		flags_codes |= 1 << 15
	}
	flags_codes |= uint16(mh.op_code&0xF) << 11
	if mh.aa {
		flags_codes |= 1 << 10
	}
	if mh.tc {
		flags_codes |= 1 << 9
	}
	if mh.rd {
		flags_codes |= 1 << 8
	}
	if mh.ra {
		flags_codes |= 1 << 7
	}
	flags_codes |= uint16(mh.r_code & 0xF)

	header_raw = binary.BigEndian.AppendUint16(header_raw, flags_codes)
	header_raw = binary.BigEndian.AppendUint16(header_raw, mh.qd_count)
	header_raw = binary.BigEndian.AppendUint16(header_raw, mh.an_count)
	header_raw = binary.BigEndian.AppendUint16(header_raw, mh.ns_count)
	header_raw = binary.BigEndian.AppendUint16(header_raw, mh.ar_count)
	return header_raw
}
func printHeader(mh *MessageHeader) {
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

type Question struct {
	Qname  []string
	Qtype  uint16
	Qclass uint16
}

func decodeName(name_raw []byte) ([]string, int) {
	name := make([]string, 0)
	index := 0
	for name_raw[index] != 0 && index < len(name_raw) {
		label_size := name_raw[index]
		var label_builder strings.Builder
		for range label_size {
			index++
			label_builder.WriteRune(rune(name_raw[index]))
		}
		name = append(name, label_builder.String())
		index++
	}

	return name, index + 1
}

func encodeName(name []string) []byte {
	name_raw := make([]byte, 0, 512)
	for _, label := range name {
		label_size := byte(len(label))
		label_raw := []byte(label)
		name_raw = append(name_raw, label_size)
		name_raw = append(name_raw, label_raw...)
	}
	name_raw = append(name_raw, byte(0))
	return name_raw
}

func decodeQuestion(q_raw []byte) (Question, int) {
	q := Question{}
	Qname, index := decodeName(q_raw)
	q.Qname = Qname
	q.Qtype = binary.BigEndian.Uint16(q_raw[index : index+2])
	index += 2
	q.Qclass = binary.BigEndian.Uint16(q_raw[index : index+2])
	index += 2
	return q, index
}

func encodeQuestion(q *Question) []byte {
	buffer := make([]byte, 0, 4096)
	buffer = append(buffer, encodeName(q.Qname)...)
	buffer = binary.BigEndian.AppendUint16(buffer, q.Qtype)
	buffer = binary.BigEndian.AppendUint16(buffer, q.Qclass)
	return buffer
}

func printQuestion(q *Question) {
	fmt.Println("Question: ")
	fmt.Print("Name: ")
	for _, label := range q.Qname {
		fmt.Printf("%s.", label)
	}
	fmt.Print("\n")
	fmt.Println("Type: ", q.Qtype)
	fmt.Println("Class: ", q.Qclass)
}

type RecordDataIpv4 struct {
	netip.Addr
}
type RecordDataIpv6 struct {
	netip.Addr
}
type RecordDataText struct {
	string
}

type RecordData struct {
	ipv4   *RecordDataIpv4
	ipv6   *RecordDataIpv6
	text   *RecordDataText
	buffer []byte
}

type ResourceRecord struct {
	Name     []string
	Type     RRType
	Class    uint16
	TTL      uint32
	RDLength uint16
	RDData   RecordData
}

func decodeResourceRecord(rr_raw []byte) (ResourceRecord, int) {
	rr := ResourceRecord{}
	Qname, index := decodeName(rr_raw)
	rr.Name = Qname
	rr.Type = RRType(binary.BigEndian.Uint16(rr_raw[index : index+2]))
	index += 2
	rr.Class = binary.BigEndian.Uint16(rr_raw[index : index+2])
	index += 2
	rr.TTL = binary.BigEndian.Uint32(rr_raw[index : index+4])
	index += 4
	rr.RDLength = binary.BigEndian.Uint16(rr_raw[index : index+2])
	index += 2
	switch rr.Type {
	case ADDRESS_IPV4:
		{
			rr.RDData = RecordData{
				ipv4: &RecordDataIpv4{},
			}
			//TODO: check if you need to flip them before assigning them in the array
			rr.RDData.ipv4.Addr = netip.AddrFrom4(
				[4]byte(rr_raw[index : index+int(rr.RDLength)]))
		}
	case ADDRESS_IPV6:
		{
			rr.RDData = RecordData{
				ipv6: &RecordDataIpv6{},
			}
			rr.RDData.ipv6.Addr = netip.AddrFrom16(
				[16]byte(rr_raw[index : index+int(rr.RDLength)]))
		}
	case TEXT:
		{
			rr.RDData = RecordData{
				text: &RecordDataText{},
			}
			var builder strings.Builder
			for i := range rr.RDLength {
				builder.WriteRune(rune(rr_raw[index+int(i)]))
			}
			rr.RDData.text.string = builder.String()
		}
	default:
		{
			rr.RDData.buffer = make([]byte, 0)
			rr.RDData.buffer = append(rr.RDData.buffer, rr_raw[index:index+int(rr.RDLength)]...)
			fmt.Println("Un supported valuetype")
		}
	}

	index += int(rr.RDLength)
	return rr, index
}

func encodeResourceRecord(rr *ResourceRecord) []byte {
	rr_raw := make([]byte, 0, 4096)
	rr_raw = append(rr_raw, encodeName(rr.Name)...)
	rr_raw = binary.BigEndian.AppendUint16(rr_raw, uint16(rr.Type))
	rr_raw = binary.BigEndian.AppendUint16(rr_raw, rr.Class)
	rr_raw = binary.BigEndian.AppendUint32(rr_raw, rr.TTL)
	rr_raw = binary.BigEndian.AppendUint16(rr_raw, rr.RDLength)
	switch rr.Type {
	case ADDRESS_IPV4:
		{
			temp := rr.RDData.ipv4.As4()
			rr_raw = append(rr_raw, temp[:]...)
		}
	case ADDRESS_IPV6:
		{
			temp := rr.RDData.ipv6.As16()
			rr_raw = append(rr_raw, temp[:]...)
		}
	case TEXT:
		{
			rr_raw = append(rr_raw, []byte(rr.RDData.text.string)...)
		}
	default:
		{
			fmt.Println("Unsupported Record type write fake data")
			rr_raw = append(rr_raw, rr.RDData.buffer...)
		}
	}
	return rr_raw
}

func printResourceRecord(rr *ResourceRecord) {
	fmt.Println("Resource Record: ")
	fmt.Print("Name: ")
	for _, label := range rr.Name {
		fmt.Printf("%s.", label)
	}
	fmt.Print("\n")
	fmt.Println("Type: ", rr.Type)
	fmt.Println("Class: ", rr.Class)
	fmt.Println("Time to Live: ", rr.TTL)
	fmt.Println("Record Data Section Length: ", rr.RDLength)
	switch rr.Type {
	case ADDRESS_IPV4:
		{
			fmt.Println("IPV4 :", rr.RDData.ipv4.String())
		}
	case ADDRESS_IPV6:
		{
			fmt.Println("IPV6 :", rr.RDData.ipv6.String())
		}
	case TEXT:
		{
			fmt.Println("Text :", rr.RDData.text.string)
		}
	}
}

type Message struct {
	header      MessageHeader
	questions   []Question
	answers     []ResourceRecord
	authorities []ResourceRecord
	additionals []ResourceRecord
} //used to represent the shape of a dns message when sending or reseaving

func decodeMessage(message_raw []byte) Message {
	message := Message{}
	index := 0
	message.header, _ = decodeHeader(message_raw[:12])
	index = 12
	for range message.header.qd_count {
		question, size := decodeQuestion(message_raw[index:])
		message.questions = append(message.questions, question)
		index += int(size)
	}
	for range message.header.an_count {
		answer, size := decodeResourceRecord(message_raw[index:])
		message.answers = append(message.answers, answer)
		index += int(size)
	}
	for range message.header.ns_count {
		authority, size := decodeResourceRecord(message_raw[index:])
		message.authorities = append(message.authorities, authority)
		index += int(size)
	}
	for range message.header.ar_count {
		additional, size := decodeResourceRecord(message_raw[index:])
		message.additionals = append(message.additionals, additional)
		index += int(size)
	}
	return message
}
func encodeMessage(message Message) []byte {
	message_raw := make([]byte, 0, 4096)
	message_raw = append(message_raw, encodeHeader(&message.header)...)
	for _, question := range message.questions {
		message_raw = append(message_raw, encodeQuestion(&question)...)
	}
	for _, answer := range message.answers {
		message_raw = append(message_raw, encodeResourceRecord(&answer)...)
	}
	for _, authority := range message.authorities {
		message_raw = append(message_raw, encodeResourceRecord(&authority)...)
	}
	for _, additional := range message.additionals {
		message_raw = append(message_raw, encodeResourceRecord(&additional)...)
	}
	return message_raw
}
func printMessage(message Message) {
	fmt.Println("Message: ")
	printHeader(&message.header)
	for i, question := range message.questions {
		fmt.Println("Question number ", i)
		printQuestion(&question)
	}
	for i, answer := range message.answers {
		fmt.Println("Answer number ", i)
		printResourceRecord(&answer)
	}
	for i, authority := range message.authorities {
		fmt.Println("Authority number ", i)
		printResourceRecord(&authority)
	}
	for i, additional := range message.additionals {
		fmt.Println("Additional number ", i)
		printResourceRecord(&additional)
	}

}

func responde(space NameSpace, message_raw []byte) []byte {
	message := decodeMessage(message_raw)
	message.header.qr = true
	message.header.aa = true
	for _, question := range message.questions {
		answer := ResourceRecord{
			Name:  question.Qname,
			Type:  RRType(question.Qtype),
			Class: question.Qclass,
			TTL:   5 * 60 * 60,
		}
		address, err := getAddress(space, question.Qname)
		if err != nil {
			message.header.r_code = 3
			break
		} else {
			if address.Is4() {
				answer.RDData.ipv4 = &RecordDataIpv4{}
				answer.RDData.ipv4.Addr = address
				answer.RDLength = uint16(answer.RDData.ipv4.BitLen()) / 8
			} else {
				answer.RDData.ipv6 = &RecordDataIpv6{}
				answer.RDData.ipv6.Addr = address
				answer.RDLength = uint16(answer.RDData.ipv6.BitLen()) / 8
			}
			message.answers = append(message.answers, answer)
		}
	}
	printResourceRecord(&message.additionals[0])
	//BUG: the message is malformed when being encoded look into what could be the problem
	message.header.an_count = uint16(len(message.answers))
	return encodeMessage(message)
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
	space := loadNameSpace()
	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("Read error: %v\n", err)
			continue
		}

		message_raw := buffer[:n]
		response := responde(space, message_raw)
		response = slices.Clip(response)
		fmt.Println("this is the actuall response body :", response)
		_, err = conn.WriteToUDP(response, clientAddr)
	}
}
