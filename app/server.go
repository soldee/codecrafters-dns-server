package main

import (
	"fmt"
	"net"
)

type Server interface {
	Run() error
}

type DefaultServer struct {
	conn *net.UDPConn
}

func NewDefaultServer(serverConn *net.UDPConn) *DefaultServer {
	return &DefaultServer{
		conn: serverConn,
	}
}

func (server *DefaultServer) Run() error {
	fmt.Println("Running default server")

	buf := make([]byte, 512)
	for {
		size, source, err := server.conn.ReadFromUDP(buf)
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

		_, err = server.conn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
	return nil
}

func generateDnsMessageResponse(receivedMessage *DnsMessage) *DnsMessage {
	var rcode uint8 = 0
	if receivedMessage.Header.Flags.OPCODE != 0 {
		rcode = 4
	}
	fmt.Println("hi1")
	anCount := receivedMessage.Header.QDCOUNT
	answers := make([]DnsAnswer, anCount)
	for i := 0; i < int(anCount); i++ {
		answers[i] = DnsAnswer{
			NAME:     receivedMessage.Questions[i].QNAME,
			TYPE:     TYPE_A,
			CLASS:    CLASS_IN,
			TTL:      60,
			RDLENGTH: 4,
			RDATA:    []byte{8, 8, 8, 8},
		}
	}
	fmt.Println("hi2")
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
			QDCOUNT: receivedMessage.Header.QDCOUNT,
			ANCOUNT: anCount,
			NSCOUNT: 0,
			ARCOUNT: 0,
		},
		Questions: receivedMessage.Questions,
		Answers:   answers,
	}
}
