package main

import (
	"fmt"
	"net"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	udpAddr, err := net.ResolveUDPAddr("udp", "192.168.33.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer udpConn.Close()

	buf := make([]byte, 512)

	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		receivedData := buf[:size]
		fmt.Printf("Request DNS message bytes: %08b\n", receivedData)
		fmt.Printf("Request DNS message hex from %s: %v\n", source, receivedData)

		receivedMessage := deserializeMessage(receivedData)
		fmt.Printf("Request DNS message: %+v\n", receivedMessage)

		var responseMessage *DnsMessage = generateDnsMessageResponse(receivedMessage)
		fmt.Printf("Response DNS message: %+v\n", responseMessage)
		var response = responseMessage.serialize()
		fmt.Printf("Response DNS message bytes: %08b\n", response)
		fmt.Printf("Response DNS message hex: %v", response)

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}

func generateDnsMessageResponse(receivedMessage *DnsMessage) *DnsMessage {
	var rcode uint8 = 0
	if receivedMessage.Header.Flags.OPCODE != 0 {
		rcode = 4
	}

	return &DnsMessage{
		Header: DnsHeader{
			ID: receivedMessage.Header.ID,
			Flags: HeaderFlags{
				QR:     true,
				OPCODE: receivedMessage.Header.Flags.OPCODE,
				AA:     false,
				TC:     false,
				RD:     receivedMessage.Header.Flags.RD,
				RA:     false,
				Z:      0,
				RCODE:  rcode,
			},
			QDCOUNT: 1,
			ANCOUNT: 1,
			NSCOUNT: 0,
			ARCOUNT: 0,
		},
		Question: DnsQuestion{
			QNAME:  receivedMessage.Question.QNAME,
			QTYPE:  TYPE_A,
			QCLASS: CLASS_IN,
		},
		Answer: DnsAnswer{
			NAME:     receivedMessage.Question.QNAME,
			TYPE:     TYPE_A,
			CLASS:    CLASS_IN,
			TTL:      60,
			RDLENGTH: 4,
			RDATA:    []byte{8, 8, 8, 8},
		},
	}
}
