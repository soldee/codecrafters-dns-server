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

		receivedData := string(buf[:size])
		fmt.Printf("Received %d bytes from %s: %s\n", size, source, receivedData)

		var response *DnsMessage = generateDnsMessageResponse()

		fmt.Printf("%08b\n", response.serialize())
		fmt.Println(response.serialize())

		_, err = udpConn.WriteToUDP(response.serialize(), source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}

func generateDnsMessageResponse() *DnsMessage {
	return &DnsMessage{
		Header: DnsHeader{
			ID: 1234,
			Flags: HeaderFlags{
				QR:     true,
				OPCODE: 0,
				AA:     false,
				TC:     false,
				RD:     false,
				RA:     false,
				Z:      0,
				RCODE:  0,
			},
			QDCOUNT: 1,
			ANCOUNT: 0,
			NSCOUNT: 0,
			ARCOUNT: 0,
		},
		Question: DnsQuestion{
			QNAME:  "codecrafters.io",
			QTYPE:  1,
			QCLASS: 1,
		},
	}
}
